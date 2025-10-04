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
	tmpQueueType, err := beego.AppConfig.DIY("queue.type")
	if err != nil {
		return cnf, err
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

		brokerURL = fmt.Sprintf("amqp://%s:%s@%s:%d/",
			user,
			password,
			host,
			port)
	}

	runMode, _ := helper.AppRunmode()
	tmpDefaultQueue, err := beego.AppConfig.DIY("queue.default_queue")
	if err != nil {
		return cnf, err
	}
	cnf = &config.Config{
		Broker:        brokerURL,
		DefaultQueue:  tmpDefaultQueue.(string) + "_" + runMode,
		ResultBackend: brokerURL,
	}

	return
}

func (mq *MqServer) TestConnection() error {
	backend := mq.Server.GetBackend()
	_, err := backend.GetState("non-existent-task-id")
	if err != nil {
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
