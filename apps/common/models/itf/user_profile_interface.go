package itf

import (
	"time"
)

type UserProfileItf interface {
	TableName() string

	GetId() string
	GetAvatar() string
	GetCardType() int
	GetCardNum() string
	GetCardImages() string
	GetGender() int
	GetBirthDate() *time.Time
	GetConstellation() string
	GetOccupation() string
	GetCompany() string
	GetEmergencyName() string
	GetEmergencyTel() string
	GetAddress() string
	GetEmail() string
	GetSource() int
	GetValidDateBegin() *time.Time
	GetValidDateEnd() *time.Time
	GetGraduatedFrom() string
	GetSchooling() string
	GetDegreeNumber() string
	GetProfessional() string
	GetStatus() int
	GetCreatedAt() int64
	GetUpdatedAt() int64
	GetDeletedAt() *int64
	GetDeleted() int
	GetRemark() string
}
