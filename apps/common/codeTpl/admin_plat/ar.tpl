{{if .IsMultiApp}}
package models_ar

import (
	"WenBeego/apps/common/models"
	{{if .HasDeleted -}}
	_ "context"
	{{- end}}
)

type {{.ModelName}}Ar struct {
	{{- if and (.IsMultiTenant) (eq .AppModule "admin_plat") }}
	models.{{.PlatModelName}}
	{{- else if and (.IsMultiTenant) ( eq .AppModule "admin_mchnt")}}
	models.{{.MchntModelName}}
	{{- else}}
	models.{{.ModelName}}
	{{- end}}
}
{{else}}
package models_ar

import (
	"WenBeego/apps/common/global"
	"WenBeego/apps/common/models"
	"WenBeego/apps/common/models/base_model"
	"WenBeego/apps/common/dto/{{.MenuModule}}_dto"
	"fmt"
	"time"

	"gorm.io/gorm"
	{{if .HasDeleted -}}
	_ "context"
	{{- end}}
)

type {{.ModelName}}Ar struct {
	base_model.{{.ModelName}}
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
		Omit("id", {{if .HasCreateTime}}"{{.CreateTimeField}}",{{end}}  {{if .HasDeleted}}"{{.DeletedField}}",{{end}} ).
		Where("id = ?", data.Id).
		Updates(data).Error
}

func (ar *{{.ModelName}}Ar) Delete(tx *gorm.DB, id string) error {
	if id == "" {
		return fmt.Errorf("{{.ModelName}}Ar Delete(): Id 不能为空")
	}
	{{if .HasDeleted}}
	return tx.Model(&models.{{.ModelName}}{}).Where("id = ?", id).Update("{{.DeletedField}}", 1).Error
	{{else}}
	ctx := context.Background()
	return tx.Model(&models.{{.ModelName}}{}).Where("id = ?", id).Delete(ctx).Error
	{{end}}
}

func (ar *{{.ModelName}}Ar) GetList(pageSize, offset int, searchDto {{.MenuModule}}_dto.{{.ModelName}}Dto) (data []models.{{.ModelName}}, count int64, err error) {
	data = make([]models.{{.ModelName}}, 0)
	model := &models.{{.ModelName}}{}
	tableName := model.TableName()

	query := global.GetReadDb().
		Model(model){{if .HasDeleted}}.
		Where(tableName + ".{{.DeletedField}} = 0"){{end}}
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
		Select({{.ListSelectCols}}).
		Limit(pageSize).
		Offset(offset).
		{{if .HasCreateTime}}Order(tableName + ".{{.CreateTimeField}} desc").{{end}}
		Find(&data).Error
	return
}

func (ar *{{.ModelName}}Ar) GetById(id string) (data base_model.{{.ModelName}}, err error) {
	err = global.GetReadDb().
		Model(&models.{{.ModelName}}{}).
		Where("{{if .HasDeleted}}{{.DeletedField}} = 0 AND {{end}}id = ?", id).
		Take(&data).Error
	return
}
{{end}}
