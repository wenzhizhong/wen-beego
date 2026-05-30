package itf

type ApiStatisticsItf interface {
	TableName() string

	GetID() string
	GetPermsID() string
	GetURI() string
	GetPV() int
	GetUV() int
	GetDate() int64
	GetUnitId() string
	GetModuleName() string
}
