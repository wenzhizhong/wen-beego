package models_ar

import (
	"WenBeego/apps/common/global"
	"WenBeego/apps/common/models"
)

type FileAr struct {
	models.File
}

func (m *FileAr) Insert(data *models.File) error {
	return global.GetWriteDb().Create(&data).Error
}

func (m *FileAr) GetByName(name string) (models.File, error) {
	var data models.File
	err := global.GetReadDb().Where("name = ?", name).First(&data).Error
	return data, err
}
