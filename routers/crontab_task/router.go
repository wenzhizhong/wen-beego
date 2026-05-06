package crontab_task

import (
	"WenBeego/apps/common/global/constant"
	"WenBeego/apps/cron_task/controllers"
)

type CronTasks struct {
	Name     string `json:"name_en"`
	NameText string `json:"name"`
	CallBack func() `json:"-"`
}

func GetCronTasks() []CronTasks {
	TasksList := []CronTasks{}

	TasksList = append(TasksList, CronTasks{
		Name:     string(constant.CRON_BIRTH_NOTICE),
		NameText: "生日提醒",
		CallBack: func() {
			(&controllers.BirthNotice{}).Notice()
		},
	})
	return TasksList
}
