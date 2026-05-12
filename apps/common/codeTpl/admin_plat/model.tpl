package models

import "WenBeego/apps/common/models/base_model"

type {{.ModelName}} struct {
	base_model.{{.ModelName}}
}

func (m *{{.ModelName}}) TableName() string {
	return `{{.TableName}}`
}
