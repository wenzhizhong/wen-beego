package itf

import (
	"time"
)

type MenuItf interface {
	TableName() string
	GetId() string
	SetId(string)

	GetParentId() string
	SetParentId(string)

	GetUnitId() string
	SetUnitId(string)

	GetMenuType() int
	SetMenuType(int)

	GetTitle() string
	SetTitle(string)

	GetName() string
	SetName(string)

	GetPath() string
	SetPath(string)

	GetComponent() string
	SetComponent(string)

	GetRank() *int
	SetRank(*int)

	GetRedirect() string
	SetRedirect(string)

	GetIcon() string
	SetIcon(string)

	GetExtraIcon() string
	SetExtraIcon(string)

	GetEnterTransition() string
	SetEnterTransition(string)

	GetLeaveTransition() string
	SetLeaveTransition(string)

	GetActivePath() string
	SetActivePath(string)

	GetAuths() string
	SetAuths(string)

	GetFrameSrc() string
	SetFrameSrc(string)

	GetFrameLoading() bool
	SetFrameLoading(bool)

	GetKeepAlive() bool
	SetKeepAlive(bool)

	GetHiddenTag() bool
	SetHiddenTag(bool)

	GetFixedTag() bool
	SetFixedTag(bool)

	GetShowLink() bool
	SetShowLink(bool)

	GetShowParent() bool
	SetShowParent(bool)

	GetCreatedAt() time.Time
	SetCreatedAt(time.Time)

	GetUpdatedAt() time.Time
	SetUpdatedAt(time.Time)

	GetDeleted() int
	SetDeleted(int)

	GetClone() int
	SetClone(int)

	GetMenuFrom() string
	SetMenuFrom(string)
}
