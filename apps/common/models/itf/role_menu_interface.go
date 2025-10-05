package itf

type RoleMenuItf interface {
	TableName() string

	GetId() string
	GetRoleId() string
	GetMenuId() string
	GetMenuPermsId() string
}
