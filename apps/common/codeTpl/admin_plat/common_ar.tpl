package models_ar

import (
	"WenBeego/apps/common/global"
	"WenBeego/apps/common/models"
	"WenBeego/apps/common/models/base_model"
	"fmt"

	"gorm.io/gorm"
)

type {{.ModelName}}Ar struct {
	models.{{.ModelName}}
}

func (ar *{{.ModelName}}Ar) Insert(tx *gorm.DB, data *base_model.{{.ModelName}}) error {
	if data.Id == "" {
		return fmt.Errorf("{{.ModelName}}Ar Insert(): Id 不能为空")
	}
	return tx.Model(&models.{{.ModelName}}{}).Create(data).Error
}

func (ar *{{.ModelName}}Ar) Update(tx *gorm.DB, data *base_model.{{.ModelName}}) error {
	if data.Id == "" {
		return fmt.Errorf("{{.ModelName}}Ar Update(): Id 不能为空")
	}
	return tx.Model(&models.{{.ModelName}}{}).
		Select("*").
		Omit("id", "created_at", "deleted").
		Where("id = ?", data.Id).
		Updates(data).Error
}

func (ar *{{.ModelName}}Ar) Delete(tx *gorm.DB, id string) error {
	if id == "" {
		return fmt.Errorf("{{.ModelName}}Ar Delete(): Id 不能为空")
	}
	return tx.Model(&models.{{.ModelName}}{}).Where("id = ?", id).Update("deleted", 1).Error
}

func (ar *{{.ModelName}}Ar) GetList(pageSize, offset int, keyword string) (data []models.{{.ModelName}}, count int64, err error) {
	data = make([]models.{{.ModelName}}, 0)
	model := &models.{{.ModelName}}{}
	tableName := model.TableName()

	query := global.GetReadDb().
		Model(model).
		Where(tableName + ".deleted = 0")

	if keyword != "" {
		query = query.Where(tableName+".\"name\" LIKE ?", "%"+keyword+"%")
	}

	err = query.Count(&count).Error
	if err != nil || count == 0 {
		return
	}

	err = query.
		Select("{{.ListSelectCols}}").
		Limit(pageSize).
		Offset(offset).
		Order(tableName + ".created_at desc").
		Find(&data).Error
	return
}

func (ar *{{.ModelName}}Ar) GetById(id string) (data base_model.{{.ModelName}}, err error) {
	err = global.GetReadDb().
		Model(&models.{{.ModelName}}{}).
		Where("deleted = 0 AND id = ?", id).
		Take(&data).Error
	return
}
