package base_ar

import (
	"WenBeego/apps/common/global"
	_ "WenBeego/apps/common/models"
	"WenBeego/apps/common/models/base_model"
	"WenBeego/apps/common/dto_vo/{{.MenuModule}}_dto"
	{{if .HasUnitId}}"strings"{{end}}
	"fmt"
	"time"

	"gorm.io/gorm"
)

// Get{{.ModelName}}List 多用户体系通用列表查询
func Get{{.ModelName}}List[M, UserM{{if .HasUnitId}}, UnitM{{end}} interface{ TableName() string }](pageSize, offset int, searchDto {{.MenuModule}}_dto.{{.ModelName}}Dto, model M, userModel UserM{{if .HasUnitId}}, unitModel UnitM{{end}}) (data []base_model.{{.ModelName}}, count int64, err error) {
	data = make([]base_model.{{.ModelName}}, 0)
	tableName := model.TableName()
	userTableName := userModel.TableName()

	query := global.GetReadDb().
		Model(model).
		Table(tableName){{if .HasDeleted}}.
		Where(tableName + ".{{.DeletedField}} = 0"){{end}}
{{if .HasCreateUserId}}	query = query.Joins("LEFT JOIN " + userTableName + " AS creator ON creator.id = " + tableName + ".{{.CreateUserIdField}} AND creator.deleted = 0")
{{end}}{{if .HasUpdateUserId}}	query = query.Joins("LEFT JOIN " + userTableName + " AS updater ON updater.id = " + tableName + ".{{.UpdateUserIdField}} AND updater.deleted = 0")
{{end}}{{if .HasUnitId}}	unitTableName := unitModel.TableName()
	query = query.Joins("LEFT JOIN " + unitTableName + " ON " + unitTableName + ".id = " + tableName + ".{{.UnitIdField}} AND " + unitTableName + ".deleted = 0")
{{end}}
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
	{{- if .HasUnitId}}
	if searchDto.SelectUnitIds != "" {
		selectUnitIds := strings.Split(searchDto.SelectUnitIds, ",")
		query = query.Where(tableName+".unit_id in (?)", selectUnitIds)
	}
	{{- end}}

	err = query.Count(&count).Error
	if err != nil || count == 0 {
		return
	}

	err = query.
		Select({{.ListSelectCols}}{{if .HasCreateUserId}}, "creator.name AS created_by_name"{{end}}{{if .HasUpdateUserId}}, "updater.name AS updated_by_name"{{end}}{{if .HasUnitId}}, unitTableName + ".name AS unit_name"{{end}}).
		Limit(pageSize).
		Offset(offset).
		{{if .HasCreateTime}}Order(tableName + ".{{.CreateTimeField}} desc").{{end}}
		Find(&data).Error
	return
}

// Get{{.ModelName}}ById 获取单条记录
func Get{{.ModelName}}ById[M interface{ TableName() string }](id string, model M) (data base_model.{{.ModelName}}, err error) {
	err = global.GetReadDb().
		Model(model).
		Where("{{if .HasDeleted}}{{.DeletedField}} = 0 AND {{end}}id = ?", id).
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
		Omit("id", {{if .HasCreateTime}}"{{.CreateTimeField}}",{{end}}  {{if .HasDeleted}}"{{.DeletedField}}",{{end}} ).
		Where("id = ?", data.Id).
		Updates(data).Error
}

// Delete{{.ModelName}} 通用软删除
func Delete{{.ModelName}}[M interface{ TableName() string }](tx *gorm.DB, id string, model M) error {
	if id == "" {
		return fmt.Errorf("Delete{{.ModelName}}: Id 不能为空")
	}
	{{if .HasDeleted}}
	return tx.Model(model).Where("id = ?", id).Update("{{.DeletedField}}", 1).Error
	{{else}}
	ctx := context.Background()
	return tx.Model(model).Where("id = ?", id).Delete(ctx).Error
	{{end}}
}
