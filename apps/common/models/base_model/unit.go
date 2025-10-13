package base_model

type Unit struct {
	Id            string `json:"id" gorm:"type:bpchar(36);not null;primaryKey;comment:ID"`
	Pid           string `json:"pid" gorm:"type:bpchar(36);not null;default:'';comment:上级组织id"`
	Logo          string `json:"logo" gorm:"type:varchar(512);default:'';comment:logo"`
	Name          string `json:"name" gorm:"type:varchar(100);not null;comment:单位名称"`
	Code          string `json:"code" gorm:"type:varchar(100);not null;comment:组织机构代码"`
	Corporation   string `json:"corporation" gorm:"type:varchar(100);not null;comment:法人"`
	License       string `json:"license" gorm:"type:varchar(512);not null;default:'';comment:营业执照"`
	Address       string `json:"address" gorm:"type:varchar(255);default:'';comment:地址"`
	Status        int    `json:"status" gorm:"type:int4;not null;default:0;comment:0未审核，1审核通过，2审核不通过，3禁用"`
	Deleted       int    `json:"deleted" gorm:"type:int4;default:0;comment:是否删除：0否1是"`
	CreatedAt     int64  `json:"created_at" gorm:"autoCreateTime;comment:创建时间"`
	UpdatedAt     int64  `json:"updated_at" gorm:"autoCreateTime;comment:更新时间"`
	DeletedAt     *int64 `json:"deleted_at" gorm:"comment:删除时间"`
	DefaultUnitId string `json:"default_unit_id" gorm:"-"`
}

var UNIT_STATUS_UNREVIEWED = 0
var UNIT_STATUS_PASSED = 1
var UNIT_STATUS_UNPASSED = 2
var UNIT_STATUS_DISABLED = 3
var UNIT_STATUS_MAP = map[int]string{
	UNIT_STATUS_UNREVIEWED: "未审核",
	UNIT_STATUS_PASSED:     "审核通过",
	UNIT_STATUS_UNPASSED:   "审核不通过",
	UNIT_STATUS_DISABLED:   "已禁用",
}
