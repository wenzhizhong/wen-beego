package models_ar

import (
	"WenBeego/apps/common/global"
	"WenBeego/apps/common/models"
	"errors"
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

// GetLinkById
func (m *FileAr) GetLinkById(ids []string) (data []models.File, err error) {
	if len(ids) == 0 {
		return data, errors.New("id不能为空")
	}
	err = global.GetReadDb().Model(&models.File{}).Select("id,real_name,path,file_md5").Where("id in ?", ids).Find(&data).Error
	return data, err
}
