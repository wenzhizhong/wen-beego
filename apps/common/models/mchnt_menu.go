package models

import (
	"WenBeego/apps/common/models/base_model"
	"WenBeego/apps/common/models/itf"
	"time"
)

var _ itf.MenuItf = (*MchntMenu)(nil)

type MchntMenu struct {
	base_model.UnitMenu
}

func (m *MchntMenu) TableName() string {
	return `mchnt_menu`
}
func (m *MchntMenu) GetId() string {
	return m.Id
}
func (m *MchntMenu) GetParentId() string {
	return m.ParentId
}
func (m *MchntMenu) GetUnitId() string {
	return m.UnitId
}
func (m *MchntMenu) GetMenuType() int {
	return m.MenuType
}
func (m *MchntMenu) GetTitle() string {
	return m.Title
}
func (m *MchntMenu) GetName() string {
	return m.Name
}
func (m *MchntMenu) GetPath() string {
	return m.Path
}
func (m *MchntMenu) GetComponent() string {
	return m.Component
}
func (m *MchntMenu) GetRank() *int {
	return m.Rank
}
func (m *MchntMenu) GetRedirect() string {
	return m.Redirect
}
func (m *MchntMenu) GetIcon() string {
	return m.Icon
}
func (m *MchntMenu) GetExtraIcon() string {
	return m.ExtraIcon
}
func (m *MchntMenu) GetEnterTransition() string {
	return m.EnterTransition
}
func (m *MchntMenu) GetLeaveTransition() string {
	return m.LeaveTransition
}
func (m *MchntMenu) GetActivePath() string {
	return m.ActivePath
}
func (m *MchntMenu) GetAuths() string {
	return m.Auths
}
func (m *MchntMenu) GetFrameSrc() string {
	return m.FrameSrc
}
func (m *MchntMenu) GetFrameLoading() bool {
	return m.FrameLoading
}
func (m *MchntMenu) GetKeepAlive() bool {
	return m.KeepAlive
}
func (m *MchntMenu) GetHiddenTag() bool {
	return m.HiddenTag
}
func (m *MchntMenu) GetFixedTag() bool {
	return m.FixedTag
}
func (m *MchntMenu) GetShowLink() bool {
	return m.ShowLink
}
func (m *MchntMenu) GetShowParent() bool {
	return m.ShowParent
}
func (m *MchntMenu) GetCreatedAt() time.Time {
	return m.CreatedAt
}
func (m *MchntMenu) GetUpdatedAt() time.Time {
	return m.UpdatedAt
}
func (m *MchntMenu) GetDeleted() int {
	return m.Deleted
}
func (m *MchntMenu) GetClone() int {
	return m.Clone
}
