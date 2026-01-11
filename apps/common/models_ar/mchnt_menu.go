package models_ar

import (
	"WenBeego/apps/common/global"
	"WenBeego/apps/common/models"
	"errors"
)

type MchntMenuAr struct {
	models.MchntMenu
}

func (a *MchntMenuAr) GetById(id string) (models.MchntMenu, error) {
	data := models.MchntMenu{}
	if id == "" {
		return data, errors.New("id不能为空")
	}
	result := global.GetReadDb().Where("id = ?", id).Take(&data)
	return data, result.Error
}
