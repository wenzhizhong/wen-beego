package models

import "WenBeego/apps/common/models/base_model"

type PlatCronLog struct {
	base_model.UnitCronLog
}

func (m *PlatCronLog) TableName() string {
	return "plat_crontab_log"
}
