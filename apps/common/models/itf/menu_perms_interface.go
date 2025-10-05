package itf

import (
	"time"
)

type MenuPermsItf interface {
	TableName() string

	GetId() string
	GetMenuId() string
	GetType() int16
	GetName() string
	GetPermission() string
	GetUri() string
	GetMethod() int16
	GetDeleted() int16
	GetCreatedAt() time.Time
	GetUpdatedAt() time.Time
}
