package models_ar

import (
	"WenBeego/apps/common/global"
	"WenBeego/apps/common/models"
	"WenBeego/apps/common/models/base_model"
	"WenBeego/apps/common/dto/{{.MenuModule}}_dto"
	"fmt"
	"time"

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

func (ar *{{.ModelName}}Ar) GetList(pageSize, offset int, searchDto {{.MenuModule}}_dto.{{.ModelName}}Dto) (data []models.{{.ModelName}}, count int64, err error) {
	data = make([]models.{{.ModelName}}, 0)
	model := &models.{{.ModelName}}{}
	tableName := model.TableName()

	query := global.GetReadDb().
		Model(model).
		Where(tableName + ".deleted = 0")
{{range .Columns}}{{if not (eq .Name "id")}}{{if not (isHasDeletedFields .Name)}}{{if not (eq .FormType "editor")}}	{{if or (eq .GoType "time.Time") (eq .GoType "*time.Time")}}	if !searchDto.{{.GoFieldName}}Start.IsZero() {
			query = query.Where(tableName+".{{.Name}} >= ?", searchDto.{{.GoFieldName}}Start)
		}
		if !searchDto.{{.GoFieldName}}End.IsZero() {
			endVal := searchDto.{{.GoFieldName}}End
			{{if eq .Type "date"}}endVal = endVal.AddDate(0, 0, 1){{else}}endVal = endVal.Add(time.Second){{end}}
			query = query.Where(tableName+".{{.Name}} < ?", endVal)
		}
	{{else if or (eq .GoType "int") (eq .GoType "int64") (eq .GoType "float32") (eq .GoType "float64")}}	if searchDto.{{.GoFieldName}} != 0 { query = query.Where(tableName+".{{.Name}} = ?", searchDto.{{.GoFieldName}}) }
	{{else if eq .GoType "bool"}}	if searchDto.{{.GoFieldName}} { query = query.Where(tableName+".{{.Name}} = ?", searchDto.{{.GoFieldName}}) }
	{{else}}	if searchDto.{{.GoFieldName}} != "" { query = query.Where(tableName+".{{.Name}} = ?", searchDto.{{.GoFieldName}}) }
	{{end}}
{{end}}{{end}}{{end}}{{end}}

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
