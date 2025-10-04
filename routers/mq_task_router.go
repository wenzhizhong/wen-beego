package routers

import (
	"WenBeego/apps/mq_task/controllers"
)

type MqTasks struct {
	Name     string
	CallBack interface{}
}

func GetMqTasks() []MqTasks {
	TasksList := []MqTasks{}

	TasksList = append(TasksList, MqTasks{Name: "ApiLog.ActionSaveToDb", CallBack: (*controllers.ApiLog).ActionSaveToDb})
	return TasksList
}
