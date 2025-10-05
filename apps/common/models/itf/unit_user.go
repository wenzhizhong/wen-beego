package itf

type UnitUserItf interface {
	TableName() string

	GetId() string
	GetUnitId() string
	GetIsDefault() int
	GetDefaultUnitId() string
	GetUserId() string
	GetDeleted() int
	GetPhone() string
	GetName() string
}
