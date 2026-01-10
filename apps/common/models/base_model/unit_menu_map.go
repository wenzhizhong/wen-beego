package base_model

type UnitMenuMap struct {
	Id        string `gorm:"column:id;type:bpchar(36);primaryKey" json:"id"`
	UnitId    string `gorm:"column:unit_id;type:bpchar(36);not null" json:"unit_id"`
	MenuId    string `gorm:"column:menu_id;type:bpchar(36);not null" json:"menu_id"`
	UpdatedAt int64  `gorm:"column:updated_at;type:int8;default:0" json:"updated_at"`
	Deleted   int    `gorm:"column:deleted;type:int4;default:0" json:"deleted"`
}
