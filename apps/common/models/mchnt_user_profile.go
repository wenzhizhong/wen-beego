package models

import (
	"WenBeego/apps/common/models/base_model"
	"WenBeego/apps/common/models/itf"
	"time"
)

var _ itf.UserProfileItf = (*MchntUserProfile)(nil)

type MchntUserProfile struct {
	base_model.UnitUserProfile
}

func (m *MchntUserProfile) TableName() string {
	return `mchnt_user_profile`
}

func (m *MchntUserProfile) GetId() string {
	return m.Id
}
func (m *MchntUserProfile) GetAvatar() string {
	return m.Avatar
}
func (m *MchntUserProfile) GetCardType() int {
	return m.CardType
}
func (m *MchntUserProfile) GetCardNum() string {
	return m.CardNum
}
func (m *MchntUserProfile) GetCardImages() string {
	return m.CardImages
}
func (m *MchntUserProfile) GetGender() int {
	return m.Gender
}
func (m *MchntUserProfile) GetBirthDate() *time.Time {
	return m.BirthDate
}
func (m *MchntUserProfile) GetConstellation() string {
	return m.Constellation
}
func (m *MchntUserProfile) GetOccupation() string {
	return m.Occupation
}
func (m *MchntUserProfile) GetCompany() string {
	return m.Company
}
func (m *MchntUserProfile) GetEmergencyName() string {
	return m.EmergencyName
}
func (m *MchntUserProfile) GetEmergencyTel() string {
	return m.EmergencyTel
}
func (m *MchntUserProfile) GetAddress() string {
	return m.Address
}
func (m *MchntUserProfile) GetEmail() string {
	return m.Email
}
func (m *MchntUserProfile) GetSource() string {
	return m.Source
}
func (m *MchntUserProfile) GetValidDateBegin() *time.Time {
	return m.ValidDateBegin
}
func (m *MchntUserProfile) GetValidDateEnd() *time.Time {
	return m.ValidDateEnd
}
func (m *MchntUserProfile) GetSchooling() string {
	return m.Schooling
}
func (m *MchntUserProfile) GetDegreeNumber() string {
	return m.DegreeNumber
}
func (m *MchntUserProfile) GetLearnProfessional() string {
	return m.LearnProfessional
}
func (m *MchntUserProfile) GetProfessional() string {
	return m.Professional
}
func (m *MchntUserProfile) GetStatus() int {
	return m.Status
}
func (m *MchntUserProfile) GetCreatedAt() int64 {
	return m.CreatedAt
}
func (m *MchntUserProfile) GetUpdatedAt() int64 {
	return m.UpdatedAt
}
func (m *MchntUserProfile) GetDeletedAt() *int64 {
	return m.DeletedAt
}
func (m *MchntUserProfile) GetDeleted() int {
	return m.Deleted
}
