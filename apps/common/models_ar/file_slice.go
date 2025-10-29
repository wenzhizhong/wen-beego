package models_ar

import (
	"WenBeego/apps/common/global"
	"WenBeego/apps/common/models"
	"errors"
)

type FileSliceAr struct {
	models.FileSlice
}

func (m *FileSliceAr) GetById(id string) (data models.FileSlice, err error) {
	if id == "" {
		return data, errors.New("id不能为空")
	}
	err = global.GetReadDb().Model(&models.FileSlice{}).Where("id = ?", id).Take(&data).Error
	return data, err
}

func (m *FileSliceAr) GetListByFileMd5(fileMd5 string, createBy string, sliceTotal ...int64) (data []models.FileSlice, err error) {
	query := global.GetReadDb().
		Model(&models.FileSlice{}).
		Where("file_md5 = ?", fileMd5).
		Where("create_by = ?", createBy)
	if len(sliceTotal) > 0 {
		query = query.Where("slice_index = ?", sliceTotal[0])
	}
	err = query.Order("slice_index").Find(&data).Error
	return data, err
}

func (m *FileSliceAr) Insert(data *models.FileSlice) error {
	return global.GetWriteDb().Create(&data).Error
}

func (m *FileSliceAr) Delete(fileMd5 string) error {
	return global.GetWriteDb().Where("file_md5 = ?", fileMd5).Delete(&models.FileSlice{}).Error
}
