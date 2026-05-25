package {{.MenuModule}}_vo

import (
	"WenBeego/apps/common/models/base_model"
)

type {{.ModelName}}Vo struct {
	base_model.{{.ModelName}}
{{if .HasCreateUserId}}	CreatedByName string `gorm:"->" json:"created_by_name"`
{{end}}{{if .HasUpdateUserId}}	UpdatedByName string `gorm:"->" json:"updated_by_name"`
{{end}}{{if .HasUnitId}}	UnitName      string `gorm:"->" json:"unit_name"`
{{end}}
}
