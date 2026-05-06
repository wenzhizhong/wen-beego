package mq

import (
	"WenBeego/apps/common/global"
	"WenBeego/apps/common/global/constant"
	"WenBeego/apps/common/helper"
	"fmt"

	"WenBeego/apps/common/thirdPkg/rewrite/RichardKnop/machinery/v1/backends/result"
	"WenBeego/apps/common/thirdPkg/rewrite/RichardKnop/machinery/v1/tasks"
)

type MqClient struct {
	MqServer
}

func (mq *MqClient) Init() {
	if global.MqClient != nil {
		mq.Server = global.MqClient
		return
	}
	mqServer, err := mq.NewMq()
	if err != nil {
		global.Log.Error("MqServer.NewMq() error:", err)
		panic(err)
	}

	err = mqServer.TestConnection()
	if err != nil {
		panic(err)
	}

	// 添加到全局变量
	global.MqClient = mqServer.Server
	global.MqClientDlx = mqServer.ServerDlx
}

// 获取基础任务签名
func (mq *MqClient) GetNewBaseSignature() (*tasks.Signature, error) {
	uuid, err := helper.GetUuid()
	if err != nil {
		return nil, err
	}
	runMode, _ := helper.AppRunmode()
	signature := &tasks.Signature{}
	signature.UUID = fmt.Sprintf("task_%v", uuid)
	signature.RetryCount = 3
	signature.IgnoreWhenTaskNotRegistered = runMode == "prod"
	return signature, nil
}

/*
*
  - 发送任务
  - @param taskName string 任务名称
  - @param args []tasks.Arg 任务参数，
*/
func (mq *MqClient) SendTask(taskName constant.MqNameType, args []tasks.Arg) (asyncResult *result.AsyncResult, err error) {
	signature, err := mq.GetNewBaseSignature()
	if err != nil {
		return nil, err
	}
	signature.Name = string(taskName)
	signature.Args = args

	asyncResult, err = global.MqClient.SendTask(signature)
	return
}
