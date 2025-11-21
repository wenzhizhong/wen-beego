package base_model

type UnitUserDept struct {
	Id      string `json:"id" gorm:"type:bpchar(36);not null;primaryKey;comment:ID"`
	UserId  string `json:"user_id" gorm:"type:bpchar(36);not null;comment:组织单位用户id"`
	DeptId  string `json:"dept_id" gorm:"type:bpchar(36);not null;comment:组织单位部门id"`
	Deleted int    `json:"deleted" gorm:"type:int4;default:0;comment:是否删除：0否1是"`
}
