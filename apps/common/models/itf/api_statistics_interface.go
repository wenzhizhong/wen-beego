package itf

type ApiStatisticsItf interface {
	TableName() string

	GetID() string
	GetPermsID() string
	GetURI() string
	GetPV() int64
	GetUV() int64
	GetDate() int64
}
