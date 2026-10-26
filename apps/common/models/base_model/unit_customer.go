package base_model

type UnitCustomer struct {
	Id      string `json:"id" gorm:"type:bpchar(36);not null;primaryKey;comment:ID"`
	UnitId  string `json:"unit_id" gorm:"type:bpchar(36);not null;comment:当前用户所属组织单位id"`
	UserId  string `json:"user_id" gorm:"type:bpchar(36);not null;comment:客户ID"`
	Status  int    `json:"status" gorm:"type:int4;not null;default:1;comment:状态"`
	Deleted int    `json:"deleted" gorm:"type:int2;not null;default:0;comment:是否删除：0否1是"`
}
