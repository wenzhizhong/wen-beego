package itf

import "time"

type UnitCronItf interface {
	TableName() string

	GetID() string
	GetUnitID() string
	GetName() string
	GetNameEn() string
	GetGroup() string
	GetCronExpr() string
	GetStatus() int
	GetCreatedBy() string
	GetCreatedAt() *time.Time
	GetUpdatedBy() *string
	GetUpdatedAt() *time.Time
	GetDeleted() int
}
