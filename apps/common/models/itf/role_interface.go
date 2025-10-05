package itf

import (
	"time"
)

type RoleItf interface {
	TableName() string

	GetId() string
	GetUnitId() string
	GetRoleName() string
	GetRoleSort() int
	GetStatus() int
	GetDeleted() int
	GetCreatedBy() string
	GetCreatedAt() time.Time
	GetUpdatedBy() string
	GetUpdatedAt() time.Time
	GetRemark() string
}
