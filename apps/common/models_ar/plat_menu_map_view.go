package models_ar

import (
	"WenBeego/apps/common/global"
	"WenBeego/apps/common/models"
	"errors"
)

type PlatMenuMapViewAr struct {
	models.PlatMenuView
}

func (a *PlatMenuMapViewAr) GetById(id string) (models.PlatMenu, error) {
	data := models.PlatMenu{}
	if id == "" {
		return data, errors.New("id不能为空")
	}
	result := global.GetReadDb().Where("id = ?", id).Take(&data)
	return data, result.Error
}
