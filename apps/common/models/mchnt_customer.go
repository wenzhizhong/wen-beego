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
func (m *MchntCustomer) GetCustomerId() string {
	return m.CustomerId
}
func (m *MchntCustomer) SetCustomerId(customerId string) {
	m.CustomerId = customerId
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
