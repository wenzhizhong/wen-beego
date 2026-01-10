package itf

type MenuMapItf interface {
	TableName() string
	GetId() string
	GetUnitId() string
	GetMenuId() string
	GetUpdatedAt() int64
	GetDeleted() int
}
