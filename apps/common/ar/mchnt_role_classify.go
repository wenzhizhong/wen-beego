package ar

import (
	"WenBeego/apps/common/global"
	"WenBeego/apps/common/models"
	"errors"
)

type MchntRoleClassifyAr struct {
	models.MchntRoleClassify
}

func (a *MchntRoleClassifyAr) GetById(id string) (models.MchntRoleClassify, error) {
	data := models.MchntRoleClassify{}
	if id == "" {
		return data, errors.New("id不能为空")
	}
	result := global.GetReadDb().Where("id = ?", id).Take(&data)
	return data, result.Error
}
