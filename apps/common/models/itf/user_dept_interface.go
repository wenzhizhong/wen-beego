package itf

type UserDeptItf interface {
	TableName() string

	GetId() string
	GetUserId() string
	GetDeptId() string
	GetDeleted() int
}
