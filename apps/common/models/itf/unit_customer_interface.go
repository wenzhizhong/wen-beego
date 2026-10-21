package itf

type UnitCustomerItf interface {
	TableName() string

	GetId() string
	SetId(string)

	GetUnitId() string
	SetUnitId(string)

	GetCustomerId() string
	SetCustomerId(string)

	GetStatus() int
	SetStatus(int)

	GetDeleted() int
	SetDeleted(int)
}
