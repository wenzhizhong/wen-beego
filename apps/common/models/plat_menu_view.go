package models

import (
	"WenBeego/apps/common/models/base_model"
	"WenBeego/apps/common/models/itf"
	"time"
)

var _ itf.MenuItf = (*PlatMenuView)(nil)

type PlatMenuView struct {
	base_model.UnitMenu
	MenuFrom string `json:"menuFrom" gorm:"->"`
}

func (m *PlatMenuView) TableName() string {
	return `plat_menu_view`
}

func (m *PlatMenuView) GetId() string {
	return m.Id
}
func (m *PlatMenuView) GetParentId() string {
	return m.ParentId
}
func (m *PlatMenuView) GetUnitId() string {
	return m.UnitId
}
func (m *PlatMenuView) GetMenuType() int {
	return m.MenuType
}
func (m *PlatMenuView) GetTitle() string {
	return m.Title
}
func (m *PlatMenuView) GetName() string {
	return m.Name
}
func (m *PlatMenuView) GetPath() string {
	return m.Path
}
func (m *PlatMenuView) GetComponent() string {
	return m.Component
}
func (m *PlatMenuView) GetRank() *int {
	return m.Rank
}
func (m *PlatMenuView) GetRedirect() string {
	return m.Redirect
}
func (m *PlatMenuView) GetIcon() string {
	return m.Icon
}
func (m *PlatMenuView) GetExtraIcon() string {
	return m.ExtraIcon
}
func (m *PlatMenuView) GetEnterTransition() string {
	return m.EnterTransition
}
func (m *PlatMenuView) GetLeaveTransition() string {
	return m.LeaveTransition
}
func (m *PlatMenuView) GetActivePath() string {
	return m.ActivePath
}
func (m *PlatMenuView) GetAuths() string {
	return m.Auths
}
func (m *PlatMenuView) GetFrameSrc() string {
	return m.FrameSrc
}
func (m *PlatMenuView) GetFrameLoading() bool {
	return m.FrameLoading
}
func (m *PlatMenuView) GetKeepAlive() bool {
	return m.KeepAlive
}
func (m *PlatMenuView) GetHiddenTag() bool {
	return m.HiddenTag
}
func (m *PlatMenuView) GetFixedTag() bool {
	return m.FixedTag
}
func (m *PlatMenuView) GetShowLink() bool {
	return m.ShowLink
}
func (m *PlatMenuView) GetShowParent() bool {
	return m.ShowParent
}
func (m *PlatMenuView) GetCreatedAt() time.Time {
	return m.CreatedAt
}
func (m *PlatMenuView) GetUpdatedAt() time.Time {
	return m.UpdatedAt
}
func (m *PlatMenuView) GetDeleted() int {
	return m.Deleted
}
func (m *PlatMenuView) GetClone() int {
	return m.Clone
}
