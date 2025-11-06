package itf

type DeptItf interface {
	TableName() string

	GetId() string
	GetPid() string
	GetUnitId() string
	GetName() string
	GetPrincipalId() string
	GetPrincipal() string
	GetPhone() string
	GetEmail() string
	GetSort() int
	GetStatus() int
	GetDeleted() int
	GetUpdatedAt() int64
	GetDeletedAt() *int64
	GetRemark() string
}
