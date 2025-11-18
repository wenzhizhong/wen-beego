package models

import (
	"WenBeego/apps/common/models/base_model"
	"WenBeego/apps/common/models/itf"
	"time"
)

var _ itf.MenuItf = (*PlatMenu)(nil)

type PlatMenu struct {
	base_model.UnitMenu
}

func (m *PlatMenu) TableName() string {
	return `plat_menu`
}

func (m *PlatMenu) GetId() string {
	return m.Id
}
func (m *PlatMenu) GetParentId() string {
	return m.ParentId
}
func (m *PlatMenu) GetUnitId() string {
	return m.UnitId
}
func (m *PlatMenu) GetMenuType() int {
	return m.MenuType
}
func (m *PlatMenu) GetTitle() string {
	return m.Title
}
func (m *PlatMenu) GetName() string {
	return m.Name
}
func (m *PlatMenu) GetPath() string {
	return m.Path
}
func (m *PlatMenu) GetComponent() string {
	return m.Component
}
func (m *PlatMenu) GetRank() *int {
	return m.Rank
}
func (m *PlatMenu) GetRedirect() string {
	return m.Redirect
}
func (m *PlatMenu) GetIcon() string {
	return m.Icon
}
func (m *PlatMenu) GetExtraIcon() string {
	return m.ExtraIcon
}
func (m *PlatMenu) GetEnterTransition() string {
	return m.EnterTransition
}
func (m *PlatMenu) GetLeaveTransition() string {
	return m.LeaveTransition
}
func (m *PlatMenu) GetActivePath() string {
	return m.ActivePath
}
func (m *PlatMenu) GetAuths() string {
	return m.Auths
}
func (m *PlatMenu) GetFrameSrc() string {
	return m.FrameSrc
}
func (m *PlatMenu) GetFrameLoading() bool {
	return m.FrameLoading
}
func (m *PlatMenu) GetKeepAlive() bool {
	return m.KeepAlive
}
func (m *PlatMenu) GetHiddenTag() bool {
	return m.HiddenTag
}
func (m *PlatMenu) GetFixedTag() bool {
	return m.FixedTag
}
func (m *PlatMenu) GetShowLink() bool {
	return m.ShowLink
}
func (m *PlatMenu) GetShowParent() bool {
	return m.ShowParent
}
func (m *PlatMenu) GetCreatedAt() time.Time {
	return m.CreatedAt
}
func (m *PlatMenu) GetUpdatedAt() time.Time {
	return m.UpdatedAt
}
func (m *PlatMenu) GetDeleted() int {
	return m.Deleted
}
func (m *PlatMenu) GetClone() int {
	return m.Clone
}
