package ar

import (
	"WenBeego/apps/common/global"
	"WenBeego/apps/common/models"
	"errors"
)

type MchntAr struct {
	models.Mchnt
}

func (a *MchntAr) GetById(id string) (models.Mchnt, error) {
	data := models.Mchnt{}
	if id == "" {
		return data, errors.New("id不能为空")
	}
	result := global.GetReadDb().Where("id = ?", id).Take(&data)
	return data, result.Error
}
