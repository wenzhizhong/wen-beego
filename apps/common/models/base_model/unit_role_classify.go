package base_model

type UnitRoleClassify struct {
	Id      string `json:"id" gorm:"type:bpchar(36);not null;primaryKey;comment:ID"`
	RoleId  string `json:"role_id" gorm:"type:bpchar(36);not null;comment:角色id"`
	UnitId  string `json:"unit_id" gorm:"type:bpchar(36);not null;comment:plat/mchut id"`
	Name    string `json:"name" gorm:"type:varchar(100);not null;comment:角色分类"`
	Deleted int    `json:"deleted" gorm:"type:int4;default:0;comment:是否删除"`
}
