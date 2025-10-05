package itf

import (
	"time"
)

type MenuItf interface {
	TableName() string

	GetId() string
	GetUnitId() string
	GetIcon() string
	GetName() string
	GetApiPath() string
	GetPagePath() string
	GetType() int
	GetPid() string
	GetAllPid() string
	GetCreatedAt() time.Time
	GetWeight() int
	GetVisible() int
	GetDeleted() int
}
