package common

import (
	"crypto/tls"
	"fmt"
	"strings"

	amqp "github.com/rabbitmq/amqp091-go"
)

// AMQPConnector ...
type AMQPConnector struct{}

// Connect opens a connection to RabbitMQ, declares an exchange, opens a channel,
// declares and binds the queue and enables publish notifications
func (ac *AMQPConnector) Connect(urls string, urlSeparator string, tlsConfig *tls.Config, exchange, exchangeType, queueName string, queueDurable, queueDelete bool, queueBindingKey string, exchangeDeclareArgs, queueDeclareArgs, queueBindingArgs amqp.Table) (*amqp.Connection, *amqp.Channel, amqp.Queue, <-chan amqp.Confirmation, <-chan *amqp.Error, error) {
	urlsList := []string{urls}
	if urlSeparator != "" {
		urlsList = strings.Split(urls, urlSeparator)
	}

	var conn *amqp.Connection
	var channel *amqp.Channel
	var err error

	for _, url := range urlsList {
		// Connect to server
		conn, channel, err = ac.Open(url, tlsConfig)
		if err != nil {
			continue
		} else {
			break
		}
	}

	if err != nil {
		return nil, nil, amqp.Queue{}, nil, nil, err
	}

	if exchange != "" {
		// Declare an exchange
		if err = channel.ExchangeDeclare(
			exchange,            // name of the exchange
			exchangeType,        // type
			true,                // durable
			false,               // delete when complete
			false,               // internal
			false,               // noWait
			exchangeDeclareArgs, // arguments
		); err != nil {
			return conn, channel, amqp.Queue{}, nil, nil, fmt.Errorf("Exchange declare error: %s", err)
		}
	}
	if queueDeclareArgs != nil && queueDeclareArgs["x-dead-letter-exchange"] != nil {
		dlxName := queueDeclareArgs["x-dead-letter-exchange"].(string)

		if dlxName != exchange {
			// 1. 声明死信交换机
			err = channel.ExchangeDeclare(
				dlxName, "direct", true, false, false, false, nil,
			)
			if err != nil {
				return nil, nil, amqp.Queue{}, nil, nil, fmt.Errorf("Declare DLX [%s]: %s", dlxName, err)
			}

		// 2. 声明死信队列（最终目的地，绑定独立key，仅由taskFailed显式发布）
		dlqFinalName := queueName + ".dlx"
		_, err = channel.QueueDeclare(
			dlqFinalName, true, false, false, false, nil,
		)
		if err != nil {
			return nil, nil, amqp.Queue{}, nil, nil, fmt.Errorf("Declare DLQ [%s]: %s", dlqFinalName, err)
		}
		err = channel.QueueBind(
			dlqFinalName, dlqFinalName, dlxName, false, nil,
		)
		if err != nil {
			return nil, nil, amqp.Queue{}, nil, nil, fmt.Errorf("Bind DLQ [%s] to DLX [%s]: %s", dlqFinalName, dlxName, err)
		}

			// 3. 声明重试队列（带TTL，过期后回到主交换机）
			retryQueueName := queueName + ".retry"
			retryDelayMs := int32(1000)
			if retryDelay, ok := queueDeclareArgs["x-retry-delay-ms"]; ok {
				switch v := retryDelay.(type) {
				case int32:
					retryDelayMs = v
				case int:
					retryDelayMs = int32(v)
				case int64:
					retryDelayMs = int32(v)
				}
			}
			retryArgs := amqp.Table{
				"x-dead-letter-exchange":    exchange,
				"x-dead-letter-routing-key": queueBindingKey,
				"x-message-ttl":             retryDelayMs,
			}
			_, err = channel.QueueDeclare(
				retryQueueName, true, false, false, false, retryArgs,
			)
			if err != nil {
				return nil, nil, amqp.Queue{}, nil, nil, fmt.Errorf("Declare RetryQueue [%s]: %s", retryQueueName, err)
			}
			err = channel.QueueBind(
				retryQueueName, queueBindingKey, dlxName, false, nil,
			)
			if err != nil {
				return nil, nil, amqp.Queue{}, nil, nil, fmt.Errorf("Bind RetryQueue [%s] to DLX [%s]: %s", retryQueueName, dlxName, err)
			}
		}
	}

	var queue amqp.Queue
	if queueName != "" {
		if queueDeclareArgs == nil {
			queueDeclareArgs = amqp.Table{}
		}
		// Declare a queue
		queue, err = channel.QueueDeclare(
			queueName,        // name
			queueDurable,     // durable
			queueDelete,      // delete when unused
			false,            // exclusive
			false,            // no-wait
			queueDeclareArgs, // arguments
		)
		if err != nil {
			return conn, channel, amqp.Queue{}, nil, nil, fmt.Errorf("Queue declare error: %s", err)
		}

		// Bind the queue
		if err = channel.QueueBind(
			queue.Name,       // name of the queue
			queueBindingKey,  // binding key
			exchange,         // source exchange
			false,            // noWait
			queueBindingArgs, // arguments
		); err != nil {
			return conn, channel, queue, nil, nil, fmt.Errorf("Queue bind error: %s", err)
		}
	}

	// Enable publish confirmations
	if err = channel.Confirm(false); err != nil {
		return conn, channel, queue, nil, nil, fmt.Errorf("Channel could not be put into confirm mode: %s", err)
	}

	return conn, channel, queue, channel.NotifyPublish(make(chan amqp.Confirmation, 1)), conn.NotifyClose(make(chan *amqp.Error, 1)), nil
}

// DeleteQueue deletes a queue by name
func (ac *AMQPConnector) DeleteQueue(channel *amqp.Channel, queueName string) error {
	// First return value is number of messages removed
	_, err := channel.QueueDelete(
		queueName, // name
		false,     // ifUnused
		false,     // ifEmpty
		false,     // noWait
	)

	return err
}

// InspectQueue provides information about a specific queue
func (*AMQPConnector) InspectQueue(channel *amqp.Channel, queueName string) (*amqp.Queue, error) {
	queueState, err := channel.QueueInspect(queueName)
	if err != nil {
		return nil, fmt.Errorf("Queue inspect error: %s", err)
	}

	return &queueState, nil
}

// Open new RabbitMQ connection
func (ac *AMQPConnector) Open(url string, tlsConfig *tls.Config) (*amqp.Connection, *amqp.Channel, error) {
	// Connect
	// From amqp docs: DialTLS will use the provided tls.Config when it encounters an amqps:// scheme
	// and will dial a plain connection when it encounters an amqp:// scheme.
	conn, err := amqp.DialTLS(url, tlsConfig)
	if err != nil {
		return nil, nil, fmt.Errorf("Dial error: %s", err)
	}

	// Open a channel
	channel, err := conn.Channel()
	if err != nil {
		return nil, nil, fmt.Errorf("Open channel error: %s", err)
	}

	return conn, channel, nil
}

// Close connection
func (ac *AMQPConnector) Close(channel *amqp.Channel, conn *amqp.Connection) error {
	if channel != nil {
		if err := channel.Close(); err != nil {
			return fmt.Errorf("Close channel error: %s", err)
		}
	}

	if conn != nil {
		if err := conn.Close(); err != nil {
			return fmt.Errorf("Close connection error: %s", err)
		}
	}

	return nil
}
