package models

import "WenBeego/apps/common/models/base_model"

type {{.MchntModelName}} struct {
	base_model.{{.ModelName}}
}

func (m *{{.MchntModelName}}) TableName() string {
	return `{{.MchntTableName}}`
}
