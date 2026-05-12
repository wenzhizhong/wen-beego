package base_ar

import (
	"WenBeego/apps/common/global"
	"WenBeego/apps/common/models"
	"WenBeego/apps/common/models/base_model"
	"fmt"

	"gorm.io/gorm"
)

// Get{{.ModelName}}List 多用户体系通用列表查询
func Get{{.ModelName}}List[M interface{ TableName() string }](pageSize, offset int, keyword string, model M) (data []models.{{.ModelName}}, count int64, err error) {
	data = make([]models.{{.ModelName}}, 0)
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
		Select(tableName + ".*").
		Limit(pageSize).
		Offset(offset).
		Order(tableName + ".created_at desc").
		Find(&data).Error
	return
}

// Get{{.ModelName}}ById 获取单条记录
func Get{{.ModelName}}ById[M interface{ TableName() string }](id string, model M) (data base_model.{{.ModelName}}, err error) {
	err = global.GetReadDb().
		Model(model).
		Where("deleted = 0 AND id = ?", id).
		Take(&data).Error
	return
}

// Insert{{.ModelName}} 通用新增
func Insert{{.ModelName}}[M interface{ TableName() string }](tx *gorm.DB, data *base_model.{{.ModelName}}, model M) error {
	return tx.Model(model).Create(data).Error
}

// Update{{.ModelName}} 通用更新
func Update{{.ModelName}}[M interface{ TableName() string }](tx *gorm.DB, data *base_model.{{.ModelName}}, model M) error {
	return tx.Model(model).
		Select("*").
		Omit("id", "created_at", "deleted").
		Where("id = ?", data.Id).
		Updates(data).Error
}

// Delete{{.ModelName}} 通用软删除
func Delete{{.ModelName}}[M interface{ TableName() string }](tx *gorm.DB, id string, model M) error {
	if id == "" {
		return fmt.Errorf("Delete{{.ModelName}}: Id 不能为空")
	}
	return tx.Model(model).Where("id = ?", id).Update("deleted", 1).Error
}
