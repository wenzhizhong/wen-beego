package models

import (
	"WenBeego/apps/common/models/base_model"
	"WenBeego/apps/common/models/itf"
)

type MchntCustomer struct {
	base_model.UnitCustomer
}

var _ itf.UnitCustomerItf = (*MchntCustomer)(nil)

func (m *MchntCustomer) TableName() string {
	return `mchnt_customer`
}

func (m *MchntCustomer) GetId() string {
	return m.Id
}
func (m *MchntCustomer) SetId(id string) {
	m.Id = id
}
func (m *MchntCustomer) GetUnitId() string {
	return m.UnitId
}
func (m *MchntCustomer) SetUnitId(unitId string) {
	m.UnitId = unitId
}
func (m *MchntCustomer) GetUserId() string {
	return m.UserId
}
func (m *MchntCustomer) SetUserId(customerId string) {
	m.UserId = customerId
}
func (m *MchntCustomer) GetStatus() int {
	return m.Status
}
func (m *MchntCustomer) SetStatus(status int) {
	m.Status = status
}
func (m *MchntCustomer) GetDeleted() int {
	return m.Deleted
}
func (m *MchntCustomer) SetDeleted(deleted int) {
	m.Deleted = deleted
}
