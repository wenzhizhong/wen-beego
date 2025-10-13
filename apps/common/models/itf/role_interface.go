package itf

type RoleItf interface {
	TableName() string

	GetId() string
	GetUnitId() string
	GetRoleName() string
	GetRoleSort() int
	GetStatus() int
	GetDeleted() int
	GetCreatedBy() string
	GetCreatedAt() int64
	GetUpdatedBy() string
	GetUpdatedAt() int64
	GetRemark() string
}
