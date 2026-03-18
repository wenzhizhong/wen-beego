package models_ar

import (
	"WenBeego/apps/common/global"
	"WenBeego/apps/common/models"
	"errors"
)

type ConfigAr struct {
	models.Config
}

func (a *ConfigAr) GetById(id string) (models.Config, error) {
	data := models.Config{}
	if id == "" {
		return data, errors.New("id不能为空")
	}
	result := global.GetReadDb().Where("id = ?", id).Take(&data)
	return data, result.Error
}

func (a *ConfigAr) GetByNames(names []string) ([]models.Config, error) {
	data := make([]models.Config, 0)
	if len(names) == 0 {
		return data, errors.New("name不能为空")
	}
	result := global.GetReadDb().Where("name in ?", names).Where("deleted=0").Find(&data)
	return data, result.Error
}
