package base_model

import "time"

type {{.ModelName}} struct {
	Id         string     ` + "`" + `gorm:"column:id;type:bpchar(36);primaryKey;comment:ID" json:"id"` + "`" + `
{{range .Columns}}	{{.GoFieldName}} {{.GoType}} ` + "`" + `gorm:"column:{{.Name}};type:{{.PgType}}{{if .Required}};not null{{end}}{{if .DefVal}};default:{{.DefVal}}{{end}};comment:{{.Comment}}" json:"{{.JsonName}}"` + "`" + `
{{end}}	CreatedAt  *time.Time ` + "`" + `gorm:"column:created_at;type:timestamptz;comment:创建时间" json:"created_at"` + "`" + `
	UpdatedAt  *time.Time ` + "`" + `gorm:"column:updated_at;type:timestamptz;comment:更新时间" json:"updated_at"` + "`" + `
	Deleted    int        ` + "`" + `gorm:"column:deleted;type:int2;default:0;comment:是否删除" json:"deleted"` + "`" + `
}
