package base_model

{{$isTimeField := false }}
{{- range .Columns -}}	
{{- if contains .GoType "time" -}}
	{{- $isTimeField = true -}}
{{- end -}}
{{- end -}}
{{if $isTimeField}} import "time"{{end}}

type {{.ModelName}} struct {
{{range .Columns}}	{{.GoFieldName}} {{.GoType}} `gorm:"column:{{.Name}};type:{{.PgType}}{{if .Required}};not null{{end}}{{if .DefVal}};default:{{.DefVal}}{{end}};comment:{{.Comment}}" json:"{{.JsonName}}"`
{{end}}	
}
