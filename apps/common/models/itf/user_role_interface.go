package itf

type UserRoleItf interface {
	TableName() string

	GetId() string
	GetUserId() string
	GetRoleId() string
	GetDeleted() int
}
