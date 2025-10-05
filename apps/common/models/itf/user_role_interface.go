package itf

type UserRoleItf interface {
	TableName() string

	GetId() string
	GetUnitId() string
	GetUserId() string
	GetRoleId() string
	GetDeleted() int
}
