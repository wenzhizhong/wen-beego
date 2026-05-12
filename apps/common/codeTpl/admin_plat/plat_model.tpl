package models

import "WenBeego/apps/common/models/base_model"

type {{.PlatModelName}} struct {
	base_model.{{.ModelName}}
}

func (m *{{.PlatModelName}}) TableName() string {
	return `{{.PlatTableName}}`
}
