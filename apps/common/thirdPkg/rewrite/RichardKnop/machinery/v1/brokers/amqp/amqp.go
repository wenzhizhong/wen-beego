package amqp

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"WenBeego/apps/common/thirdPkg/rewrite/RichardKnop/machinery/v1/brokers/errs"
	"WenBeego/apps/common/thirdPkg/rewrite/RichardKnop/machinery/v1/brokers/iface"
	"WenBeego/apps/common/thirdPkg/rewrite/RichardKnop/machinery/v1/common"
	"WenBeego/apps/common/thirdPkg/rewrite/RichardKnop/machinery/v1/config"
	"WenBeego/apps/common/thirdPkg/rewrite/RichardKnop/machinery/v1/log"
	"WenBeego/apps/common/thirdPkg/rewrite/RichardKnop/machinery/v1/tasks"

	"github.com/pkg/errors"
	amqp "github.com/rabbitmq/amqp091-go"
)

type AMQPConnection struct {
	queueName    string
	connection   *amqp.Connection
	channel      *amqp.Channel
	queue        amqp.Queue
	confirmation <-chan amqp.Confirmation
	errorchan    <-chan *amqp.Error
	cleanup      chan struct{}
}

// Broker represents an AMQP broker
type Broker struct {
	common.Broker
	common.AMQPConnector
	processingWG sync.WaitGroup

	connections      map[string]*AMQPConnection
	connectionsMutex sync.RWMutex
}

// New creates new Broker instance
func New(cnf *config.Config) iface.Broker {
	return &Broker{
		Broker:        common.NewBroker(cnf),
		AMQPConnector: common.AMQPConnector{},
		connections:   make(map[string]*AMQPConnection),
	}
}

// StartConsuming enters a loop and waits for incoming messages
func (b *Broker) StartConsuming(consumerTag string, concurrency int, taskProcessor iface.TaskProcessor) (bool, error) {
	b.Broker.StartConsuming(consumerTag, concurrency, taskProcessor)

	queueName := taskProcessor.CustomQueue()
	if queueName == "" {
		queueName = b.GetConfig().DefaultQueue
	}

	// 1. 初始化队列参数
	queueDeclareArgs := amqp.Table{}
	if b.GetConfig().AMQP != nil && b.GetConfig().AMQP.QueueDeclareArgs != nil {
		for k, v := range b.GetConfig().AMQP.QueueDeclareArgs {
			queueDeclareArgs[k] = v
		}
	}

	// 2. 为【主队列】配置死信交换机 -> 路由到重试队列
	if b.GetConfig().AMQP != nil && b.GetConfig().AMQP.DeadLetterExchange != "" {
		queueDeclareArgs["x-dead-letter-exchange"] = b.GetConfig().AMQP.DeadLetterExchange
		queueDeclareArgs["x-dead-letter-routing-key"] = b.GetConfig().AMQP.BindingKey
		if b.GetConfig().AMQP.RetryDelayMs > 0 {
			queueDeclareArgs["x-retry-delay-ms"] = b.GetConfig().AMQP.RetryDelayMs
		}
	}

	// 3. 连接到【主队列】
	conn, channel, queue, _, amqpCloseChan, err := b.Connect(
		b.GetConfig().Broker,
		b.GetConfig().MultipleBrokerSeparator,
		b.GetConfig().TLSConfig,
		b.GetConfig().AMQP.Exchange,
		b.GetConfig().AMQP.ExchangeType,
		queueName,
		true,  // queue durable
		false, // queue delete when unused
		b.GetConfig().AMQP.BindingKey,
		nil,
		queueDeclareArgs,
		amqp.Table(b.GetConfig().AMQP.QueueBindingArgs),
	)
	if err != nil {
		b.GetRetryFunc()(b.GetRetryStopChan())
		return b.GetRetry(), err
	}
	defer b.Close(channel, conn)

	// 4. 设置预取
	if err = channel.Qos(
		b.GetConfig().AMQP.PrefetchCount,
		0,     // prefetch size
		false, // global
	); err != nil {
		return b.GetRetry(), fmt.Errorf("Channel qos error: %s", err)
	}

	// 5. 消费【主队列】
	deliveries, err := channel.Consume(
		queue.Name,
		consumerTag,
		false, // auto-ack: FALSE
		false, // exclusive
		false, // no-local
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		return b.GetRetry(), fmt.Errorf("Queue consume error: %s", err)
	}

	log.INFO.Printf("[*] Consuming from queue: %s", queue.Name)

	if err := b.consume(deliveries, concurrency, taskProcessor, amqpCloseChan); err != nil {
		return b.GetRetry(), err
	}

	b.processingWG.Wait()
	return b.GetRetry(), nil
}

// StopConsuming quits the loop
func (b *Broker) StopConsuming() {
	b.Broker.StopConsuming()
	b.processingWG.Wait()
}

func (b *Broker) StartDLQConsuming(dlqQueue string, handler iface.DLQHandler) error {
	conn, channel, queue, _, amqpCloseChan, err := b.Connect(
		b.GetConfig().Broker,
		b.GetConfig().MultipleBrokerSeparator,
		b.GetConfig().TLSConfig,
		b.GetConfig().AMQP.Exchange,
		b.GetConfig().AMQP.ExchangeType,
		dlqQueue,
		true,
		false,
		dlqQueue,
		nil,
		nil,
		nil,
	)
	if err != nil {
		return fmt.Errorf("DLQ connect error: %s", err)
	}
	defer b.Close(channel, conn)

	if err = channel.Qos(1, 0, false); err != nil {
		return fmt.Errorf("DLQ Qos error: %s", err)
	}

	deliveries, err := channel.Consume(
		queue.Name, "dlq_consumer", false, false, false, false, nil,
	)
	if err != nil {
		return fmt.Errorf("DLQ consume error: %s", err)
	}

	log.INFO.Printf("[*] DLQ consuming from: %s", queue.Name)

	for {
		select {
		case amqpErr := <-amqpCloseChan:
			return amqpErr
		case d, ok := <-deliveries:
			if !ok {
				return nil
			}
			sig := new(tasks.Signature)
			if err := json.Unmarshal(d.Body, sig); err != nil {
				log.ERROR.Printf("DLQ unmarshal error: %s", err)
				d.Nack(false, true)
				continue
			}
			_ = handler(sig)
			d.Ack(false)
		case <-b.GetStopChan():
			return nil
		}
	}
}

// GetOrOpenConnection returns a connection for a queue
// ... existing code ...
func (b *Broker) GetOrOpenConnection(queueName string, queueBindingKey string, exchangeDeclareArgs, queueDeclareArgs, queueBindingArgs amqp.Table) (*AMQPConnection, error) {
	b.connectionsMutex.Lock()
	defer b.connectionsMutex.Unlock()

	if conn, ok := b.connections[queueName]; ok {
		return conn, nil
	}

	// 1. 初始化队列参数
	finalQueueArgs := amqp.Table{}
	if queueDeclareArgs != nil {
		for k, v := range queueDeclareArgs {
			finalQueueArgs[k] = v
		}
	}

	// 2. 为【主队列】配置死信交换机 (强制覆盖) -> 路由到重试队列
	if queueName == b.GetConfig().DefaultQueue &&
		b.GetConfig().AMQP != nil &&
		b.GetConfig().AMQP.DeadLetterExchange != "" {

		finalQueueArgs["x-dead-letter-exchange"] = b.GetConfig().AMQP.DeadLetterExchange
		finalQueueArgs["x-dead-letter-routing-key"] = b.GetConfig().AMQP.BindingKey
		if b.GetConfig().AMQP.RetryDelayMs > 0 {
			finalQueueArgs["x-retry-delay-ms"] = b.GetConfig().AMQP.RetryDelayMs
		}
	}

	// 3. 连接队列
	conn := &AMQPConnection{
		queueName: queueName,
		cleanup:   make(chan struct{}),
	}
	var err error
	conn.connection, conn.channel, conn.queue, conn.confirmation, conn.errorchan, err = b.Connect(
		b.GetConfig().Broker,
		b.GetConfig().MultipleBrokerSeparator,
		b.GetConfig().TLSConfig,
		b.GetConfig().AMQP.Exchange,
		b.GetConfig().AMQP.ExchangeType,
		queueName,
		true,
		false,
		queueBindingKey,
		exchangeDeclareArgs,
		finalQueueArgs,
		queueBindingArgs,
	)
	if err != nil {
		return nil, errors.Wrapf(err, "Failed to connect to queue %s", queueName)
	}

	// 4. 自动重连
	go func() {
		select {
		case err := <-conn.errorchan:
			log.INFO.Printf("Error on queue %s: %v. Reconnecting...", queueName, err)
			b.connectionsMutex.Lock()
			delete(b.connections, queueName)
			b.connectionsMutex.Unlock()
			b.GetOrOpenConnection(queueName, queueBindingKey, exchangeDeclareArgs, queueDeclareArgs, queueBindingArgs)
		case <-conn.cleanup:
			return
		}
	}()

	b.connections[queueName] = conn
	return conn, nil
}

// ... existing code ...

// ... existing code ...

// CloseConnections closes all connections
func (b *Broker) CloseConnections() error {
	b.connectionsMutex.Lock()
	defer b.connectionsMutex.Unlock()

	for _, conn := range b.connections {
		if err := b.Close(conn.channel, conn.connection); err != nil {
			log.ERROR.Print("Failed to close channel")
		}
		close(conn.cleanup)
	}
	b.connections = make(map[string]*AMQPConnection)
	return nil
}

// Publish places a new message on the queue
func (b *Broker) Publish(ctx context.Context, signature *tasks.Signature) error {
	b.AdjustRoutingKey(signature)

	msg, err := json.Marshal(signature)
	if err != nil {
		return fmt.Errorf("JSON marshal error: %s", err)
	}

	// 1. 处理延迟任务
	if signature.ETA != nil {
		now := time.Now().UTC()
		if signature.ETA.After(now) {
			delayMs := int64(signature.ETA.Sub(now) / time.Millisecond)
			return b.delay(signature, delayMs)
		}
	}

	// 2. 准备发布到【主队列】的参数
	queue := b.GetConfig().DefaultQueue
	bindingKey := b.GetConfig().AMQP.BindingKey
	if b.isDirectExchange() {
		queue = signature.RoutingKey
		bindingKey = signature.RoutingKey
	}

	// 【终极强制注入】无论上层传什么，这里必须构造出带 DLX 的参数
	publishQueueArgs := amqp.Table{
		"x-dead-letter-exchange":    b.GetConfig().AMQP.DeadLetterExchange,
		"x-dead-letter-routing-key": b.GetConfig().AMQP.BindingKey,
	}
	if b.GetConfig().AMQP.RetryDelayMs > 0 {
		publishQueueArgs["x-retry-delay-ms"] = b.GetConfig().AMQP.RetryDelayMs
	}

	// 合并其他可能存在的参数
	if b.GetConfig().AMQP.QueueDeclareArgs != nil {
		for k, v := range b.GetConfig().AMQP.QueueDeclareArgs {
			publishQueueArgs[k] = v
		}
	}

	conn, err := b.GetOrOpenConnection(
		queue,
		bindingKey,
		nil,
		publishQueueArgs, // 传入我们刚构造的、绝对正确的参数
		amqp.Table(b.GetConfig().AMQP.QueueBindingArgs),
	)
	if err != nil {
		return errors.Wrapf(err, "Failed to get connection for queue %s", queue)
	}

	channel := conn.channel
	confirmsChan := conn.confirmation

	// 3. 发布消息
	if err := channel.Publish(
		b.GetConfig().AMQP.Exchange,
		signature.RoutingKey,
		false,
		false,
		amqp.Publishing{
			Headers:      amqp.Table(signature.Headers),
			ContentType:  "application/json",
			Body:         msg,
			Priority:     signature.Priority,
			DeliveryMode: amqp.Persistent,
		},
	); err != nil {
		return errors.Wrap(err, "Failed to publish task")
	}

	confirmed := <-confirmsChan
	if !confirmed.Ack {
		return fmt.Errorf("Failed delivery of tag: %d", confirmed.DeliveryTag)
	}

	return nil
}

// consume manages worker pool
func (b *Broker) consume(deliveries <-chan amqp.Delivery, concurrency int, taskProcessor iface.TaskProcessor, amqpCloseChan <-chan *amqp.Error) error {
	pool := make(chan struct{}, concurrency)
	go func() {
		for i := 0; i < concurrency; i++ {
			pool <- struct{}{}
		}
	}()

	errorsChan := make(chan error, 1)

	for {
		select {
		case amqpErr := <-amqpCloseChan:
			return amqpErr
		case err := <-errorsChan:
			return err
		case d := <-deliveries:
			if concurrency > 0 {
				<-pool
			}

			b.processingWG.Add(1)
			go func(delivery amqp.Delivery) {
				defer func() {
					if r := recover(); r != nil {
						log.ERROR.Printf("Panic: %v", r)
						delivery.Nack(false, false) // 进入死信队列
					}
					b.processingWG.Done()
					if concurrency > 0 {
						pool <- struct{}{}
					}
				}()

				// 处理消息
				if err := b.consumeOne(delivery, taskProcessor, true); err != nil {
					log.ERROR.Printf("Task error: %v", err)
				}
			}(d)
		case <-b.GetStopChan():
			return nil
		}
	}
}

// consumeOne processes a single message
func (b *Broker) consumeOne(delivery amqp.Delivery, taskProcessor iface.TaskProcessor, ack bool) error {
	if len(delivery.Body) == 0 {
		delivery.Nack(true, false)
		return errors.New("empty message")
	}

	signature := new(tasks.Signature)
	if err := json.NewDecoder(bytes.NewReader(delivery.Body)).Decode(signature); err != nil {
		log.ERROR.Printf("Unmarshal error: %s", err)
		delivery.Nack(false, false) // 进入死信队列
		return errs.NewErrCouldNotUnmarshalTaskSignature(delivery.Body, err)
	}

	if delivery.Headers != nil {
		signature.Headers = tasks.Headers(delivery.Headers)
	}

	// 兜底保护：若重试次数已超过配置上限，发布到DLQ后Ack
	if signature.RetryCount > 0 {
		queueName := b.GetConfig().DefaultQueue
		if retryAttempts := _getXDeathRetryCount(delivery.Headers, queueName); retryAttempts > signature.RetryCount {
			log.WARNING.Printf("Task %s exceeded retry limit (%d > %d), moving to DLQ", signature.UUID, retryAttempts, signature.RetryCount)
			if err := b.PublishToDLQ(context.Background(), signature); err != nil {
				log.ERROR.Printf("Failed to publish %s to DLQ in safety net: %v", signature.UUID, err)
			}
			delivery.Ack(false)
			return nil
		} else {
			log.DEBUG.Printf("Task %s retry safety check: attempts=%d limit=%d, passing through", signature.UUID, retryAttempts, signature.RetryCount)
		}
	}

	if !b.IsTaskRegistered(signature.Name) {
		if !signature.IgnoreWhenTaskNotRegistered {
			delivery.Nack(false, true) // 重新入队
		} else {
			delivery.Ack(false)
		}
		return nil
	}

	log.DEBUG.Printf("Received: %s", delivery.Body)

	// 执行任务
	err := taskProcessor.Process(signature)
	if ack {
		if err != nil {
			// 任务失败 -> 进入【死信队列】
			log.ERROR.Printf("Task failed: %v", err)
			delivery.Nack(false, false) // 不重入队，进入死信
		} else {
			// 任务成功 -> 从队列删除
			log.DEBUG.Printf("Task succeeded")
			delivery.Ack(false)
		}
	}
	return err
}

// delay handles delayed tasks
func (b *Broker) delay(signature *tasks.Signature, delayMs int64) error {
	if delayMs <= 0 {
		return errors.New("invalid delay")
	}

	message, err := json.Marshal(signature)
	if err != nil {
		return fmt.Errorf("marshal error: %s", err)
	}

	// 1. 延迟队列名称
	queueName := b.GetConfig().AMQP.DelayedQueue
	if queueName == "" {
		queueName = fmt.Sprintf(
			"delay.%d.%s.%s",
			delayMs,
			b.GetConfig().AMQP.Exchange,
			signature.RoutingKey,
		)
	}

	// 2. 延迟队列参数
	declareQueueArgs := amqp.Table{
		// 过期后路由到【主交换机】
		"x-dead-letter-exchange":    b.GetConfig().AMQP.Exchange,
		"x-dead-letter-routing-key": signature.RoutingKey,
		"x-message-ttl":             delayMs,
		"x-expires":                 delayMs * 2,
	}

	// 3. 连接到【延迟队列】
	conn, channel, _, _, _, err := b.Connect(
		b.GetConfig().Broker,
		b.GetConfig().MultipleBrokerSeparator,
		b.GetConfig().TLSConfig,
		b.GetConfig().AMQP.Exchange,
		b.GetConfig().AMQP.ExchangeType,
		queueName,
		true,
		b.GetConfig().AMQP.AutoDelete,
		queueName,
		nil,
		declareQueueArgs,
		amqp.Table(b.GetConfig().AMQP.QueueBindingArgs),
	)
	if err != nil {
		return err
	}
	defer b.Close(channel, conn)

	// 4. 发布到延迟队列
	return channel.Publish(
		b.GetConfig().AMQP.Exchange,
		queueName,
		false,
		false,
		amqp.Publishing{
			Headers:      amqp.Table(signature.Headers),
			ContentType:  "application/json",
			Body:         message,
			DeliveryMode: amqp.Persistent,
		},
	)
}

// isDirectExchange checks if exchange type is direct
func (b *Broker) isDirectExchange() bool {
	return b.GetConfig().AMQP != nil && b.GetConfig().AMQP.ExchangeType == "direct"
}

// AdjustRoutingKey sets correct routing key
func (b *Broker) AdjustRoutingKey(s *tasks.Signature) {
	if s.RoutingKey != "" {
		return
	}
	if b.isDirectExchange() {
		s.RoutingKey = b.GetConfig().AMQP.BindingKey
	} else {
		s.RoutingKey = b.GetConfig().DefaultQueue
	}
}

func (b *Broker) PublishToDLQ(ctx context.Context, signature *tasks.Signature) error {
	dlqName := b.GetConfig().DefaultQueue + ".dlx"
	dlxName := b.GetConfig().AMQP.DeadLetterExchange
	if dlxName == "" {
		return fmt.Errorf("DeadLetterExchange not configured")
	}

	msg, err := json.Marshal(signature)
	if err != nil {
		return fmt.Errorf("JSON marshal error: %s", err)
	}

	conn, err := b.getOrOpenDLQConn(dlqName, dlxName)
	if err != nil {
		return errors.Wrapf(err, "Failed to get DLQ connection for %s", dlqName)
	}

	if err := conn.channel.Publish(
		dlxName,
		dlqName,
		false,
		false,
		amqp.Publishing{
			Headers:      amqp.Table(signature.Headers),
			ContentType:  "application/json",
			Body:         msg,
			DeliveryMode: amqp.Persistent,
		},
	); err != nil {
		return errors.Wrap(err, "Failed to publish to DLQ")
	}

	confirmed := <-conn.confirmation
	if !confirmed.Ack {
		return fmt.Errorf("Failed delivery to DLQ, tag: %d", confirmed.DeliveryTag)
	}

	return nil
}

func (b *Broker) getOrOpenDLQConn(dlqName, dlxName string) (*AMQPConnection, error) {
	connKey := "dlq:" + dlqName
	b.connectionsMutex.Lock()
	if conn, ok := b.connections[connKey]; ok {
		b.connectionsMutex.Unlock()
		return conn, nil
	}
	b.connectionsMutex.Unlock()

	connection, channel, queue, confirmsChan, errorChan, err := b.Connect(
		b.GetConfig().Broker,
		b.GetConfig().MultipleBrokerSeparator,
		b.GetConfig().TLSConfig,
		dlxName,
		"direct",
		dlqName,
		true,
		false,
		dlqName,
		nil,
		nil,
		nil,
	)
	if err != nil {
		return nil, err
	}

	conn := &AMQPConnection{
		queueName:    dlqName,
		connection:   connection,
		channel:      channel,
		queue:        queue,
		confirmation: confirmsChan,
		errorchan:    errorChan,
		cleanup:      make(chan struct{}),
	}

	b.connectionsMutex.Lock()
	b.connections[connKey] = conn
	b.connectionsMutex.Unlock()

	go func() {
		select {
		case amqpErr := <-errorChan:
			log.INFO.Printf("DLQ connection error for %s: %v. Cleaning up.", dlqName, amqpErr)
			b.connectionsMutex.Lock()
			delete(b.connections, connKey)
			b.connectionsMutex.Unlock()
		case <-conn.cleanup:
			return
		}
	}()

	return conn, nil
}

// GetPendingTasks retrieves pending tasks
func (b *Broker) GetPendingTasks(queue string) ([]*tasks.Signature, error) {
	if queue == "" {
		queue = b.GetConfig().DefaultQueue
	}

	conn, err := b.GetOrOpenConnection(
		queue,
		b.GetConfig().AMQP.BindingKey,
		nil,
		amqp.Table(b.GetConfig().AMQP.QueueDeclareArgs),
		amqp.Table(b.GetConfig().AMQP.QueueBindingArgs),
	)
	if err != nil {
		return nil, errors.Wrapf(err, "Failed to get connection for queue %s", queue)
	}

	channel := conn.channel
	queueInfo, err := channel.QueueInspect(queue)
	if err != nil {
		return nil, errors.Wrapf(err, "Failed to inspect queue %s", queue)
	}

	dumper := &sigDumper{customQueue: queue}
	for i := 0; i < queueInfo.Messages; i++ {
		d, ok, err := channel.Get(queue, false)
		if err != nil || !ok {
			break
		}
		d.Nack(false, true) // 重新入队
		b.consumeOne(d, dumper, false)
	}

	return dumper.Signatures, nil
}

// Helper for GetPendingTasks
type sigDumper struct {
	customQueue string
	Signatures  []*tasks.Signature
}

func (s *sigDumper) Process(sig *tasks.Signature) error {
	s.Signatures = append(s.Signatures, sig)
	return nil
}

func (s *sigDumper) CustomQueue() string { return s.customQueue }

func (_ *sigDumper) PreConsumeHandler() bool { return true }

func _getXDeathRetryCount(headers amqp.Table, queueName string) int {
	if headers == nil || queueName == "" {
		return 0
	}
	xDeathRaw, ok := headers["x-death"]
	if !ok {
		return 0
	}
	xDeathList, ok := xDeathRaw.([]interface{})
	if !ok {
		return 0
	}
	maxCount := 0
	for _, entry := range xDeathList {
		entryMap, ok := entry.(amqp.Table)
		if !ok {
			if m, ok2 := entry.(map[string]interface{}); ok2 {
				entryMap = amqp.Table(m)
			} else {
				continue
			}
		}
		reason, _ := entryMap["reason"].(string)
		queue, _ := entryMap["queue"].(string)
		if reason == "rejected" && queue == queueName {
			c := _getInt2(entryMap["count"])
			if c > maxCount {
				maxCount = c
			}
		}
	}
	return maxCount
}

func _getInt2(v interface{}) int {
	switch c := v.(type) {
	case int:
		return c
	case int8:
		return int(c)
	case int16:
		return int(c)
	case int32:
		return int(c)
	case int64:
		return int(c)
	case uint:
		return int(c)
	case uint8:
		return int(c)
	case uint16:
		return int(c)
	case uint32:
		return int(c)
	case uint64:
		return int(c)
	case float32:
		return int(c)
	case float64:
		return int(c)
	}
	log.ERROR.Printf("_getInt2: unknown numeric type %T value %v", v, v)
	return 0
}
