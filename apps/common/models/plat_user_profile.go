package models

import (
	"WenBeego/apps/common/models/base_model"
	"WenBeego/apps/common/models/itf"
	"time"

	"gorm.io/gorm"
)

var _ itf.UserProfileItf = (*PlatUserProfile)(nil)

type PlatUserProfile struct {
	base_model.UnitUserProfile
}

func (m *PlatUserProfile) TableName() string {
	return `plat_user_profile`
}

func (m *PlatUserProfile) GetId() string {
	return m.Id
}
func (m *PlatUserProfile) GetAvatar() string {
	return m.Avatar
}
func (m *PlatUserProfile) GetCardType() int {
	return m.CardType
}
func (m *PlatUserProfile) GetCardNum() string {
	return m.CardNum
}
func (m *PlatUserProfile) GetCardImages() string {
	return m.CardImages
}
func (m *PlatUserProfile) GetGender() int {
	return m.Gender
}
func (m *PlatUserProfile) GetBirthDate() time.Time {
	return m.BirthDate
}
func (m *PlatUserProfile) GetConstellation() string {
	return m.Constellation
}
func (m *PlatUserProfile) GetOccupation() string {
	return m.Occupation
}
func (m *PlatUserProfile) GetCompany() string {
	return m.Company
}
func (m *PlatUserProfile) GetEmergencyName() string {
	return m.EmergencyName
}
func (m *PlatUserProfile) GetEmergencyTel() string {
	return m.EmergencyTel
}
func (m *PlatUserProfile) GetAddress() string {
	return m.Address
}
func (m *PlatUserProfile) GetEMail() string {
	return m.EMail
}
func (m *PlatUserProfile) GetSource() string {
	return m.Source
}
func (m *PlatUserProfile) GetHeadimgurl() string {
	return m.Headimgurl
}
func (m *PlatUserProfile) GetValidDateBegin() time.Time {
	return m.ValidDateBegin
}
func (m *PlatUserProfile) GetValidDateEnd() time.Time {
	return m.ValidDateEnd
}
func (m *PlatUserProfile) GetSchooling() string {
	return m.Schooling
}
func (m *PlatUserProfile) GetDegreeNumber() string {
	return m.DegreeNumber
}
func (m *PlatUserProfile) GetLearnProfessional() string {
	return m.LearnProfessional
}
func (m *PlatUserProfile) GetProfessional() string {
	return m.Professional
}
func (m *PlatUserProfile) GetStatus() int {
	return m.Status
}
func (m *PlatUserProfile) GetCreatedAt() time.Time {
	return m.CreatedAt
}
func (m *PlatUserProfile) GetUpdatedAt() time.Time {
	return m.UpdatedAt
}
func (m *PlatUserProfile) GetDeletedAt() gorm.DeletedAt {
	return m.DeletedAt
}
func (m *PlatUserProfile) GetDeleted() int {
	return m.Deleted
}
