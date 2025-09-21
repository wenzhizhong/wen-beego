package ar

import (
	"WenBeego/apps/common/global"
	"WenBeego/apps/common/models"
	"errors"
)

type UserProfileAr struct {
	models.UserProfile
}

func (a *UserProfileAr) GetById(id string) (models.UserProfile, error) {
	data := models.UserProfile{}
	if id == "" {
		return data, errors.New("id不能为空")
	}
	result := global.GetReadDb().Where("id = ?", id).Take(&data)
	return data, result.Error
}
