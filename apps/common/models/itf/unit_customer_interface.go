package itf

type UnitCustomerItf interface {
	TableName() string

	GetId() string
	SetId(string)

	GetUnitId() string
	SetUnitId(string)

	GetUserId() string
	SetUserId(string)

	GetStatus() int
	SetStatus(int)

	GetDeleted() int
	SetDeleted(int)
}
