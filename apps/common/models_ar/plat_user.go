package models_ar

import (
	"WenBeego/apps/common/global"
	"WenBeego/apps/common/models"
	"errors"
)

type PlatUserAr struct {
	models.PlatUser
}

func (a *PlatUserAr) GetById(id string) (models.PlatUser, error) {
	data := models.PlatUser{}
	if id == "" {
		return data, errors.New("id不能为空")
	}
	result := global.GetReadDb().Where("id = ?", id).Take(&data)
	return data, result.Error
}
