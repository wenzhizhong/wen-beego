package controllers

import (
	"WenBeego/apps/common/dto_vo/mq_dto"
	"WenBeego/apps/common/global"
	"WenBeego/apps/mq_task/services"
	"encoding/base64"
	"encoding/json"
)

type ApiLog struct {
	servicesApiLog services.ApiLog
}

func (c *ApiLog) ActionSaveToDb(base64JsonStr string) error {
	jsonStr, err := base64.StdEncoding.DecodeString(base64JsonStr)
	if err != nil {
		return err
	}

	data := []mq_dto.ApiLogDto{}
	err = json.Unmarshal([]byte(jsonStr), &data)
	if err != nil {
		global.Log.Error("json.Unmarshal err: %v", err)
		return err
	}

	_, err = c.servicesApiLog.SaveToDb(data)
	if err != nil {
		global.Log.Error("ApiLogService.SaveToDb err: %v", err)
		return err
	}
	global.Log.Info("ApiLogService.SaveToDb success")
	return err

}
