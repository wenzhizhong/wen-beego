package models

type PlatRoleClassify struct {
	Id     string `json:"id" gorm:"type:bpchar(36);not null;primaryKey;comment:ID"`
	RoleId string `json:"role_id" gorm:"type:bpchar(36);not null;comment:角色id"`
	Name   string `json:"name" gorm:"type:varchar(100);not null;comment:角色分类"`
}

func (m *PlatRoleClassify) TableName() string {
	return `plat_role_classify`
}
