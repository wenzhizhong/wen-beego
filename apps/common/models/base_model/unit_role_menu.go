package base_model

type UnitRoleMenu struct {
	Id          string `json:"id" gorm:"type:bpchar(36);not null;primaryKey;comment:ID"`
	RoleId      string `json:"role_id" gorm:"type:bpchar(36);not null;comment:角色ID"`
	MenuId      string `json:"menu_id" gorm:"type:bpchar(36);not null;comment:菜单权限ID"`
	MenuPermsId string `json:"menu_perms_id" gorm:"type:varchar(36);not null;comment:资源关联菜单"`
}
