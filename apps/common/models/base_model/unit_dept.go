package base_model

type UnitDept struct {
	Id          string `json:"id" gorm:"type:bpchar(36);not null;primaryKey;comment:ID"`
	Pid         string `json:"pid" gorm:"type:bpchar(36);comment:上级部门id"`
	UnitId      string `json:"unit_id" gorm:"type:bpchar(36);not null;comment:组织单位id"`
	Name        string `json:"name" gorm:"type:varchar(100);not null;comment:部门名称"`
	PrincipalId string `json:"principal_id" gorm:"type:bpchar(36);comment:负责人id"`
	Principal   string `json:"principal" gorm:"->"`
	Phone       string `json:"phone" gorm:"->"`
	Email       string `json:"email" gorm:"->"`
	Sort        int    `json:"sort" gorm:"type:int4;not null;default:0;comment:排序"`
	Status      int    `json:"status" gorm:"type:int4;not null;default:0;comment:状态：0禁用1启用"`
	Deleted     int    `json:"deleted" gorm:"type:int4;not null;default:0;comment:是否删除：0否1是"`
	UpdatedAt   int64  `json:"updated_at" gorm:"type:int8;not null;comment:更新时间"`
	DeletedAt   *int64 `json:"deleted_at" gorm:"type:int8;comment:删除时间"`
	Remark      string `json:"remark" gorm:"type:varchar(512);comment:备注"`
	UnitName    string `json:"unit_name" gorm:"->"`
}

var UNIT_DEPT_STATUS_DISABLED = 0
var UNIT_DEPT_STATUS_ENABLED = 1
var UNIT_DEPT_STATUS_MAP = map[int]string{
	UNIT_DEPT_STATUS_DISABLED: "禁用",
	UNIT_DEPT_STATUS_ENABLED:  "启用",
}
