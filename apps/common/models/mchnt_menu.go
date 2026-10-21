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
func (m *MchntMenu) GetMenuFrom() string {
	return ""
}

func (m *MchntMenu) SetId(id string) {
	m.Id = id
}
func (m *MchntMenu) SetParentId(parentId string) {
	m.ParentId = parentId
}
func (m *MchntMenu) SetUnitId(unitId string) {
	m.UnitId = unitId
}
func (m *MchntMenu) SetMenuType(menuType int) {
	m.MenuType = menuType
}
func (m *MchntMenu) SetTitle(title string) {
	m.Title = title
}
func (m *MchntMenu) SetName(name string) {
	m.Name = name
}
func (m *MchntMenu) SetPath(path string) {
	m.Path = path
}
func (m *MchntMenu) SetComponent(component string) {
	m.Component = component
}
func (m *MchntMenu) SetRank(rank *int) {
	m.Rank = rank
}
func (m *MchntMenu) SetRedirect(redirect string) {
	m.Redirect = redirect
}
func (m *MchntMenu) SetIcon(icon string) {
	m.Icon = icon
}
func (m *MchntMenu) SetExtraIcon(extraIcon string) {
	m.ExtraIcon = extraIcon
}
func (m *MchntMenu) SetEnterTransition(enterTransition string) {
	m.EnterTransition = enterTransition
}
func (m *MchntMenu) SetLeaveTransition(leaveTransition string) {
	m.LeaveTransition = leaveTransition
}
func (m *MchntMenu) SetActivePath(activePath string) {
	m.ActivePath = activePath
}
func (m *MchntMenu) SetAuths(auths string) {
	m.Auths = auths
}
func (m *MchntMenu) SetFrameSrc(frameSrc string) {
	m.FrameSrc = frameSrc
}
func (m *MchntMenu) SetFrameLoading(frameLoading bool) {
	m.FrameLoading = frameLoading
}
func (m *MchntMenu) SetKeepAlive(keepAlive bool) {
	m.KeepAlive = keepAlive
}
func (m *MchntMenu) SetHiddenTag(hiddenTag bool) {
	m.HiddenTag = hiddenTag
}
func (m *MchntMenu) SetFixedTag(fixedTag bool) {
	m.FixedTag = fixedTag
}
func (m *MchntMenu) SetShowLink(showLink bool) {
	m.ShowLink = showLink
}
func (m *MchntMenu) SetShowParent(showParent bool) {
	m.ShowParent = showParent
}
func (m *MchntMenu) SetCreatedAt(createdAt time.Time) {
	m.CreatedAt = createdAt
}
func (m *MchntMenu) SetUpdatedAt(updatedAt time.Time) {
	m.UpdatedAt = updatedAt
}
func (m *MchntMenu) SetDeleted(deleted int) {
	m.Deleted = deleted
}
func (m *MchntMenu) SetClone(clone int) {
	m.Clone = clone
}
func (m *MchntMenu) SetMenuFrom(menuFrom string) {

}
