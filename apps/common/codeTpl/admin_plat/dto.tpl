package {{.MenuModule}}_dto

import (
	"WenBeego/apps/common/models/base_model"
	"time"
)

type {{.ModelName}}Dto struct {
	base_model.{{.ModelName}}
{{range .Columns}}{{if or (eq .GoType "time.Time") (eq .GoType "*time.Time")}}	{{.GoFieldName}}Start time.Time `json:"{{.Name}}Start"`
	{{.GoFieldName}}End   time.Time `json:"{{.Name}}End"`
{{end}}{{end}}
}
