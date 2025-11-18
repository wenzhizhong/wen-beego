package itf

import (
	"time"
)

type MenuItf interface {
	TableName() string
	GetId() string
	GetParentId() string
	GetUnitId() string
	GetMenuType() int
	GetTitle() string
	GetName() string
	GetPath() string
	GetComponent() string
	GetRank() *int
	GetRedirect() string
	GetIcon() string
	GetExtraIcon() string
	GetEnterTransition() string
	GetLeaveTransition() string
	GetActivePath() string
	GetAuths() string
	GetFrameSrc() string
	GetFrameLoading() bool
	GetKeepAlive() bool
	GetHiddenTag() bool
	GetFixedTag() bool
	GetShowLink() bool
	GetShowParent() bool
	GetCreatedAt() time.Time
	GetUpdatedAt() time.Time
	GetDeleted() int
	GetClone() int
}
