package base_model

type UnitApiStatistics struct {
	ID         string `gorm:"column:id;type:bpchar(36);primaryKey" json:"id"`                            // ID
	PermsID    string `gorm:"column:perms_id;type:bpchar(36);default:''" json:"perms_id"`                // menu_perms.id
	URI        string `gorm:"column:uri;type:varchar;not null;default:''" json:"uri"`                    // URI
	PV         int    `gorm:"column:pv;type:int8;not null;default:0" json:"pv"`                          // 当日PV
	UV         int    `gorm:"column:uv;type:int8;not null;default:0" json:"uv"`                          // 单日UV
	Date       int64  `gorm:"column:date;type:int8;not null" json:"date"`                                // 日期
	UnitId     string `gorm:"column:unit_id;type:bpchar(36);not null" json:"unit_id"`                    // 组织单位id
	Modulename string `gorm:"column:modulename;type:varchar(100);not null;default:''" json:"modulename"` // 模块名称
}
