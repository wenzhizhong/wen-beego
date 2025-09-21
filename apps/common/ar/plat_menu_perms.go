package ar

import (
	"WenBeego/apps/common/global"
	"WenBeego/apps/common/models"
	"errors"
)

type PlatMenuPermsAr struct {
	models.PlatMenuPerms
}

func (a *PlatMenuPermsAr) GetById(id string) (models.PlatMenuPerms, error) {
	data := models.PlatMenuPerms{}
	if id == "" {
		return data, errors.New("id不能为空")
	}
	result := global.GetReadDb().Where("id = ?", id).Take(&data)
	return data, result.Error
}
