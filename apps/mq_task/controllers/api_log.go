package controllers

import (
	"WenBeego/apps/common/dto"
	"WenBeego/apps/common/global"
	"WenBeego/apps/common/helper"
	"WenBeego/apps/mq_task/services"
	"encoding/base64"
	"encoding/json"
)

type ApiLog struct {
	servicesApiLog services.ApiLog
}

func (c *ApiLog) ActionSaveToDb(base64JsonStr string) helper.Response {
	jsonStr, err := base64.StdEncoding.DecodeString(base64JsonStr)
	if err != nil {
		panic(err)
	}
	data := []dto.ApiLogDto{}
	err = json.Unmarshal([]byte(jsonStr), &data)
	if err != nil {
		global.Log.Error("json.Unmarshal err: %v", err)
		return helper.Response{Code: 0, Message: err.Error()}
	}

	result, err := c.servicesApiLog.SaveToDb(data)
	if err != nil {
		global.Log.Error("ApiLogService.SaveToDb err: %v", err)
		return helper.Response{Code: 0, Message: err.Error()}
	}
	global.Log.Error("ApiLogService.SaveToDb success")
	return helper.Response{Code: 200, Message: "保存成功", Data: result}

}
