package models

type MchntApiStatistics struct {
	ID      string `gorm:"column:id;type:bpchar(36);primaryKey" json:"id"`             // ID
	PermsID string `gorm:"column:perms_id;type:bpchar(36);default:''" json:"perms_id"` // menu_perms.id
	URI     string `gorm:"column:uri;type:varchar;not null;default:''" json:"uri"`     // URI
	PV      int64  `gorm:"column:pv;type:int8;not null;default:0" json:"pv"`           // 当日PV
	UV      int64  `gorm:"column:uv;type:int8;not null;default:0" json:"uv"`           // 单日UV
	Date    int64  `gorm:"column:date;type:int8;not null" json:"date"`                 // 日期
}

// TableName 指定表名
func (m *MchntApiStatistics) TableName() string {
	return "mchnt_api_statistics"
}
