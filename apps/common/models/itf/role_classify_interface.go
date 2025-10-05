package itf

type RoleClassifyItf interface {
	TableName() string

	GetId() string
	GetRoleId() string
	GetUnitId() string
	GetName() string
	GetDeleted() int
}
