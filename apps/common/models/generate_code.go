package models

import "WenBeego/apps/common/models/base_model"

type GenerateCode struct {
	base_model.GenerateCode
}

func (m *GenerateCode) TableName() string {
	return `generate_code`
}
