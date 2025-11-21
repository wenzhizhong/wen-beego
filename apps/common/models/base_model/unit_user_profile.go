package base_model

import (
	"time"
)

type UnitUserProfile struct {
	Id                string     `json:"id" gorm:"not null;primaryKey;type:bpchar(36);comment:ID"`
	Avatar            string     `json:"avatar" gorm:"size:255;comment:头像"`
	CardType          int        `json:"card_type" gorm:"type:int2;comment:1大陆身份证2港澳台身份证3护照4军官证5其它"`
	CardNum           string     `json:"card_num" gorm:"size:100;comment:证件号码"`
	CardImages        string     `json:"card_images" gorm:"size:1000;comment:证件照片"`
	Gender            int        `json:"gender" gorm:"comment:性别:1男，2女"`
	BirthDate         *time.Time `json:"birth_date" gorm:"type:date;comment:出生日期"`
	Constellation     string     `json:"constellation" gorm:"size:50;comment:星座"`
	Occupation        string     `json:"occupation" gorm:"size:50;comment:职业"`
	Company           string     `json:"company" gorm:"size:500;comment:所属公司名称"`
	EmergencyName     string     `json:"emergency_name" gorm:"size:50;comment:紧急联系人姓名"`
	EmergencyTel      string     `json:"emergency_tel" gorm:"size:100;comment:紧急联系人电话"`
	Address           string     `json:"address" gorm:"size:200;comment:通讯地址"`
	Email             string     `json:"email" gorm:"size:50;comment:邮箱"`
	Source            string     `json:"source" gorm:"not null;default:'微信';comment:来源：微信,web,其它,app"`
	ValidDateBegin    *time.Time `json:"valid_date_begin" gorm:"comment:身份证有效期开始时间"`
	ValidDateEnd      *time.Time `json:"valid_date_end" gorm:"comment:身份证有效期截止时间"`
	Schooling         string     `json:"schooling" gorm:"size:100;comment:学历"`
	DegreeNumber      string     `json:"degree_number" gorm:"size:100;comment:学位编号"`
	LearnProfessional string     `json:"learn_professional" gorm:"size:100;comment:所学专业"`
	Professional      string     `json:"professional" gorm:"size:100;comment:职业"`
	Status            int        `json:"status" gorm:"not null;default:1;comment:用户行为状态：1正常，2已注销，组织单位状态：1正常，3禁用，4离职"`
	CreatedAt         int64      `json:"created_at" gorm:"autoCreateTime;comment:创建时间"`
	UpdatedAt         int64      `json:"updated_at" gorm:"autoCreateTime;comment:更新时间"`
	DeletedAt         *int64     `json:"deleted_at" gorm:"comment:删除时间"`
	Deleted           int        `json:"deleted" gorm:"not null;default:0;comment:是否删除"`
}

var UNIT_USER_PROFILE_NORMAL = 1
var UNIT_USER_PROFILE_CANCLED = 2
var UNIT_USER_PROFILE_DISABLED = 3
var UNIT_USER_PROFILE_LEAVE = 4

var UNIT_USER_PROFILE_MAP = map[int]string{
	UNIT_USER_PROFILE_NORMAL:   "正常",
	UNIT_USER_PROFILE_CANCLED:  "已注销",
	UNIT_USER_PROFILE_DISABLED: "已禁用",
	UNIT_USER_PROFILE_LEAVE:    "已离职",
}

var UNIT_GENDER_MAIL = 1
var UNIT_GENDER_FEMAIL = 2
var UNIT_GENDER_MAP = map[int]string{
	UNIT_GENDER_MAIL:   "男",
	UNIT_GENDER_FEMAIL: "女",
}

var UNIT_CARD_TYPE_1 = 1
var UNIT_CARD_TYPE_2 = 2
var UNIT_CARD_TYPE_3 = 3
var UNIT_CARD_TYPE_4 = 4
var UNIT_CARD_TYPE_5 = 5
var UNIT_CARD_TYPE_MAP = map[int]string{
	UNIT_CARD_TYPE_1: "大陆身份证",
	UNIT_CARD_TYPE_2: "港澳台身份证",
	UNIT_CARD_TYPE_3: "护照",
	UNIT_CARD_TYPE_4: "军官证",
	UNIT_CARD_TYPE_5: "其它",
}

var UNIT_USER_SOURCE_WECHAT = "微信"
var UNIT_USER_SOURCE_WEB = "web"
var UNIT_USER_SOURCE_APP = "app"
var UNIT_USER_SOURCE_OTHER = "其它"
var UNIT_USER_SOURCE_MAP = map[string]string{
	UNIT_USER_SOURCE_WECHAT: UNIT_USER_SOURCE_WECHAT,
	UNIT_USER_SOURCE_WEB:    UNIT_USER_SOURCE_WEB,
	UNIT_USER_SOURCE_APP:    UNIT_USER_SOURCE_APP,
	UNIT_USER_SOURCE_OTHER:  UNIT_USER_SOURCE_OTHER,
}
