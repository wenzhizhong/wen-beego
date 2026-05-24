package {{.MenuModule}}_dto

import (
	{{- if .IsMultiApp}}
	"WenBeego/apps/common/models/base_model"
	{{- else}}
	"WenBeego/apps/common/models"
	{{- end}}
	"time"
)

type {{.ModelName}}Dto struct {
	{{- if .IsMultiApp}}
	base_model.{{.ModelName}}
	{{- else}}
	models.{{.ModelName}}
	{{- end}}
{{range .Columns}}{{if or (eq .GoType "time.Time") (eq .GoType "*time.Time")}}	{{.GoFieldName}}Start time.Time `json:"{{.Name}}Start"`
	{{.GoFieldName}}End   time.Time `json:"{{.Name}}End"`
{{end}}{{end}}
{{if .HasUnitId}}SelectUnitIds string{{end}}
}
