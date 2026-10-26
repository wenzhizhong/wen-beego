package models

import (
	"WenBeego/apps/common/models/base_model"
	"WenBeego/apps/common/models/itf"
)

type PlatCustomer struct {
	base_model.UnitCustomer
}

var _ itf.UnitCustomerItf = (*PlatCustomer)(nil)

func (m *PlatCustomer) TableName() string {
	return `plat_customer`
}

func (m *PlatCustomer) GetId() string {
	return m.Id
}
func (m *PlatCustomer) SetId(id string) {
	m.Id = id
}
func (m *PlatCustomer) GetUnitId() string {
	return m.UnitId
}
func (m *PlatCustomer) SetUnitId(unitId string) {
	m.UnitId = unitId
}
func (m *PlatCustomer) GetUserId() string {
	return m.UserId
}
func (m *PlatCustomer) SetUserId(customerId string) {
	m.UserId = customerId
}
func (m *PlatCustomer) GetStatus() int {
	return m.Status
}
func (m *PlatCustomer) SetStatus(status int) {
	m.Status = status
}
func (m *PlatCustomer) GetDeleted() int {
	return m.Deleted
}
func (m *PlatCustomer) SetDeleted(deleted int) {
	m.Deleted = deleted
}
