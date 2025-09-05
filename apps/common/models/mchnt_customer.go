package models

type MchntCustomer struct {
	Id         string `json:"id" gorm:"type:bpchar(36);not null;primaryKey;comment:ID"`
	MchntId    string `json:"mchnt_id" gorm:"type:bpchar(36);not null;comment:商户ID"`
	CustomerId string `json:"customer_id" gorm:"type:bpchar(36);not null;comment:客户ID"`
	Status     int32  `json:"status" gorm:"type:int4;not null;default:1;comment:状态"`
	Deleted    int16  `json:"deleted" gorm:"type:int2;not null;default:0;comment:是否删除：0否1是"`
}

func (m *MchntCustomer) TableName() string {
	return `mchnt_customer`
}
