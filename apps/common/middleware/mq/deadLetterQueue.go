package mq

import (
	"WenBeego/apps/common/helper"

	beego "github.com/beego/beego/v2/server/web"
	amqp "github.com/rabbitmq/amqp091-go"
)

func getDeadLetterQueueName() (dlxName string, dlqName string, err error) {
	defQueueName := ""
	defQueueName, err = GetDefaultQueueName()

	dlxName = defQueueName + "." + "dlx"
	dlqName = defQueueName + "." + "dlq"
	return
}
func isRabbitMq() bool {
	tmpQueueType, _ := beego.AppConfig.DIY("queue.type")
	return "rabbitmq" == tmpQueueType.(string)
}

func SetupDeadLetterQueue(conn *amqp.Connection) error {
	if !isRabbitMq() {
		return nil
	}
	ch, err := conn.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()

	dlxName, dlqName, err := getDeadLetterQueueName()
	if err != nil {
		return err
	}

	// 1. 声明死信交换器
	err = ch.ExchangeDeclare(
		dlxName,  // 死信交换器名称
		"fanout", // 交换器类型
		true,     // durable: 持久化，确保RabbitMQ重启后不丢失
		false,    // auto-deleted: 不自动删除
		false,    // internal: 否
		false,    // no-wait: 否
		nil,      // arguments: 额外参数
	)
	if err != nil {
		return err
	}

	// 2. 声明死信队列
	q, err := ch.QueueDeclare(
		dlqName, // 队列名称
		true,    // durable: 持久化
		false,   // auto-delete: 不自动删除
		false,   // exclusive: 非排他
		false,   // no-wait: 否
		nil,     // arguments: 额外参数
	)
	if err != nil {
		return err
	}

	// 3. 死信队列绑定到死信交换器
	err = ch.QueueBind(
		q.Name,  // 队列名称
		"",      // routing key，fanout 交换器会忽略 routing key，传入空字符串即可
		dlxName, // 交换器名称
		false,   // no-wait
		nil,     // arguments
	)
	if err != nil {
		return err
	}

	return nil
}

func setupMainQueueWithDLX(conn *amqp.Connection) error {
	if !isRabbitMq() {
		return nil
	}
	ch, err := conn.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()

	defaultQueue, err1 := GetDefaultQueueName()
	dlxName, _, err2 := getDeadLetterQueueName()
	err = helper.Ternary(err != nil, err1, err2)
	if err != nil {
		return err
	}

	args := amqp.Table{
		"x-dead-letter-exchange": dlxName,
	}

	// 声明工作队列
	// 请确保这个队列名称和 Machinery 配置（DefaultQueue 或 BindingKey）完全一致。
	_, err = ch.QueueDeclare(
		defaultQueue, // 主队列名称
		true,         // durable
		false,        // auto-delete
		false,        // exclusive
		false,        // no-wait
		args,         // 传入死信配置
	)
	if err != nil {
		return err
	}

	return nil
}
