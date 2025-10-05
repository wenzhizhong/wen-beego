package itf

import (
	"time"
)

type UnitItf interface {
	TableName() string

	GetId() string
	GetPid() string
	GetLogo() string
	GetName() string
	GetCode() string
	GetCorporation() string
	GetLicense() string
	GetAddress() string
	GetStatus() int
	GetDeleted() int
	GetCreatedAt() time.Time
	GetUpdatedAt() time.Time
}
