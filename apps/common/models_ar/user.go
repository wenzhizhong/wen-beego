package models_ar

import (
	"WenBeego/apps/common/global"
	"WenBeego/apps/common/models"
	"errors"
)

type UserAr struct {
	models.User
}

func (a *UserAr) GetById(id string) (models.User, error) {
	data := models.User{}
	if id == "" {
		return data, errors.New("id不能为空")
	}
	result := global.GetReadDb().Where("id = ?", id).Take(&data)
	return data, result.Error
}

func (a *UserAr) GetByPhone(phone string) (models.User, error) {
	user := models.User{}
	if phone == "" {
		return user, errors.New("手机号不能为空")
	}
	result := global.GetReadDb().Where("phone = ?", phone).Take(&user)
	return user, result.Error
}
