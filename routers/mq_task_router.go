package routers

import (
	"WenBeego/apps/common/global/constant"
	"WenBeego/apps/mq_task/controllers"
)

type MqTasks struct {
	Name     string
	CallBack interface{}
}

func GetMqTasks() []MqTasks {
	TasksList := []MqTasks{}

	TasksList = append(TasksList, MqTasks{Name: string(constant.MQ_TEST_MSG), CallBack: (*controllers.Test).ActionTestMsg})
	TasksList = append(TasksList, MqTasks{Name: string(constant.MQ_API_LOG_SAVE_TO_DB), CallBack: (*controllers.ApiLog).ActionSaveToDb})

	return TasksList
}
