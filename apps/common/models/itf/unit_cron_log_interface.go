package itf

type UnitCronLogItf interface {
	TableName() string

	GetID() string
	GetCronID() string
	GetRemark() string
	GetResult() bool
	GetCreatedAt() int64
}
