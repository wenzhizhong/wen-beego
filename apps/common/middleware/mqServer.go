package middleware

import (
	"WenBeego/apps/common/global"
	"WenBeego/apps/common/helper"
	"errors"
	"fmt"
	"strings"

	"github.com/RichardKnop/machinery/v1"
	"github.com/RichardKnop/machinery/v1/config"
	beego "github.com/beego/beego/v2/server/web"
)

type MqServer struct {
	Server *machinery.Server
}

// 统一的配置创建函数
func (mq *MqServer) NewMq() (*MqServer, error) {

	cnf, err := mq.getConfig()
	if err != nil {
		return mq, err
	}

	server, err := machinery.NewServer(cnf)
	if err != nil {
		return mq, err
	}
	mq.Server = server
	return mq, nil
}

// 任务函数注册
func (mq *MqServer) RegisterTask(funcName string, funcCallBack func(...any) error) {
	mq.Server.RegisterTask(funcName, funcCallBack)
}

// 获取配置
func (mq *MqServer) getConfig() (cnf *config.Config, err error) {
	tmpQueueType, err1 := beego.AppConfig.DIY("queue.type")
	runMode, err2 := helper.AppRunmode()
	tmpDefaultQueue, err3 := beego.AppConfig.DIY("queue.default_queue")
	if err1 != nil {
		err = err1
		return cnf, err
	}
	if err2 != nil {
		err = err2
		return cnf, err
	}
	if err3 != nil {
		err = err3
		return cnf, err
	}
	defaultQueue := tmpDefaultQueue.(string) + "_" + runMode
	cnf = &config.Config{
		DefaultQueue: defaultQueue,
	}

	var brokerURL string
	queueType, _ := tmpQueueType.(string)
	switch queueType {
	case "redis":
		tmpRedis, err := beego.AppConfig.DIY("queue.redis")
		if err != nil {
			return cnf, err
		}
		tmpRedisMap, err := helper.Interface2MapInterface(tmpRedis)
		if err != nil {
			return cnf, err
		}
		if tmpRedisMap["host"] == nil || tmpRedisMap["port"] == nil || tmpRedisMap["dbNum"] == nil {
			return cnf, errors.New("queue.redis config error")
		}
		host := tmpRedisMap["host"].(string)
		port := tmpRedisMap["port"].(int)
		dbNum := tmpRedisMap["dbNum"].(int)
		password := ""
		brokerURL = fmt.Sprintf("redis://%s:%d/%d", host, port, dbNum)
		if tmpRedisMap["password"] != nil && tmpRedisMap["password"].(string) != "" {
			password = tmpRedisMap["password"].(string)
			brokerURL = fmt.Sprintf("redis://%s@%s:%d/%d", password, host, port, dbNum)
		}
		cnf.Broker = brokerURL
		cnf.ResultBackend = brokerURL
	case "rabbitmq":
		tmpRabbitMQ, err := beego.AppConfig.DIY("queue.rabbitmq")
		if err != nil {
			return cnf, err
		}
		tmpRabitMqConfig, err := helper.Interface2MapInterface(tmpRabbitMQ)
		if err != nil {
			return cnf, err
		}
		if tmpRabitMqConfig["user"] == nil || tmpRabitMqConfig["password"] == nil || tmpRabitMqConfig["host"] == nil || tmpRabitMqConfig["port"] == nil {
			return cnf, errors.New("queue.rabbitmq config error")
		}
		user := tmpRabitMqConfig["user"].(string)
		password := tmpRabitMqConfig["password"].(string)
		host := tmpRabitMqConfig["host"].(string)
		port := tmpRabitMqConfig["port"].(int)
		exchange := tmpRabitMqConfig["exchange"].(string)
		exchangeType := tmpRabitMqConfig["exchangeType"].(string)
		bindingKey := defaultQueue // tmpRabitMqConfig["bindingKey"].(string)

		brokerURL = fmt.Sprintf("amqp://%s:%s@%s:%d/",
			user,
			password,
			host,
			port)

		cnf.Broker = brokerURL
		cnf.ResultBackend = brokerURL
		cnf.AMQP = &config.AMQPConfig{
			Exchange:     exchange,
			ExchangeType: exchangeType,
			BindingKey:   bindingKey,
		}
	}

	return
}

func (mq *MqServer) TestConnection() error {
	backend := mq.Server.GetBackend()
	_, err := backend.GetState("non-existent-task-id")
	if err != nil {
		// 对于 AMQP，"No state ready" 是正常的，因为它意味着连接成功但任务不存在
		if strings.Contains(err.Error(), "No state ready") {
			global.Log.Info("✅ MQ connection test successful (no state found as expected)")
			return nil
		}
		// 忽略任务不存在的错误，关注连接错误
		if !strings.Contains(err.Error(), "not found") &&
			!strings.Contains(err.Error(), "redigo: nil returned") {
			global.Log.Error("❌ MQ backend connection test failed:", err)
			return err
		}
	}

	global.Log.Info("✅ MQ connection test successful")
	return nil
}
