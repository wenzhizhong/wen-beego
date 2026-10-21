package models

import (
	"WenBeego/apps/common/models/base_model"
	"WenBeego/apps/common/models/itf"
	"time"
)

var _ itf.MenuItf = (*PlatMenu)(nil)

type PlatMenu struct {
	base_model.UnitMenu
	MenuFrom string `json:"menuFrom" gorm:"column:menu_from;comment:系统菜单分类：admin_plat、admin_mchnt；admin_mchnt表示作为商户系统数据管理"`
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
func (m *PlatMenu) GetMenuFrom() string {
	return m.MenuFrom
}

func (m *PlatMenu) SetId(id string) {
	m.Id = id
}
func (m *PlatMenu) SetParentId(parentId string) {
	m.ParentId = parentId
}
func (m *PlatMenu) SetUnitId(unitId string) {
	m.UnitId = unitId
}
func (m *PlatMenu) SetMenuType(menuType int) {
	m.MenuType = menuType
}
func (m *PlatMenu) SetTitle(title string) {
	m.Title = title
}
func (m *PlatMenu) SetName(name string) {
	m.Name = name
}
func (m *PlatMenu) SetPath(path string) {
	m.Path = path
}
func (m *PlatMenu) SetComponent(component string) {
	m.Component = component
}
func (m *PlatMenu) SetRank(rank *int) {
	m.Rank = rank
}
func (m *PlatMenu) SetRedirect(redirect string) {
	m.Redirect = redirect
}
func (m *PlatMenu) SetIcon(icon string) {
	m.Icon = icon
}
func (m *PlatMenu) SetExtraIcon(extraIcon string) {
	m.ExtraIcon = extraIcon
}
func (m *PlatMenu) SetEnterTransition(enterTransition string) {
	m.EnterTransition = enterTransition
}
func (m *PlatMenu) SetLeaveTransition(leaveTransition string) {
	m.LeaveTransition = leaveTransition
}
func (m *PlatMenu) SetActivePath(activePath string) {
	m.ActivePath = activePath
}
func (m *PlatMenu) SetAuths(auths string) {
	m.Auths = auths
}
func (m *PlatMenu) SetFrameSrc(frameSrc string) {
	m.FrameSrc = frameSrc
}
func (m *PlatMenu) SetFrameLoading(frameLoading bool) {
	m.FrameLoading = frameLoading
}
func (m *PlatMenu) SetKeepAlive(keepAlive bool) {
	m.KeepAlive = keepAlive
}
func (m *PlatMenu) SetHiddenTag(hiddenTag bool) {
	m.HiddenTag = hiddenTag
}
func (m *PlatMenu) SetFixedTag(fixedTag bool) {
	m.FixedTag = fixedTag
}
func (m *PlatMenu) SetShowLink(showLink bool) {
	m.ShowLink = showLink
}
func (m *PlatMenu) SetShowParent(showParent bool) {
	m.ShowParent = showParent
}
func (m *PlatMenu) SetCreatedAt(createdAt time.Time) {
	m.CreatedAt = createdAt
}
func (m *PlatMenu) SetUpdatedAt(updatedAt time.Time) {
	m.UpdatedAt = updatedAt
}
func (m *PlatMenu) SetDeleted(deleted int) {
	m.Deleted = deleted
}
func (m *PlatMenu) SetClone(clone int) {
	m.Clone = clone
}
func (m *PlatMenu) SetMenuFrom(menuFrom string) {
	m.MenuFrom = menuFrom
}
