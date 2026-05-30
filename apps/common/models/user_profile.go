package models

import (
	"WenBeego/apps/common/models/base_model"
	"time"
)

type UserProfile struct {
	Id             string    `json:"id" gorm:"not null;primaryKey;type:bpchar(36);comment:ID"`
	Avatar         string    `json:"avatar" gorm:"size:255;comment:头像"`
	CardType       int       `json:"card_type" gorm:"type:int2;comment:1大陆身份证2港澳台身份证3护照4军官证5其它"`
	CardNum        string    `json:"card_num" gorm:"size:100;comment:证件号码"`
	CardImages     string    `json:"card_images" gorm:"size:1000;comment:证件照片"`
	Gender         int       `json:"gender" gorm:"comment:性别:1男，2女"`
	BirthDate      time.Time `json:"birth_date" gorm:"type:date;comment:出生日期"`
	Constellation  string    `json:"constellation" gorm:"size:50;comment:星座"`
	Occupation     string    `json:"occupation" gorm:"size:50;comment:职业"`
	Company        string    `json:"company" gorm:"size:500;comment:所属公司名称"`
	EmergencyName  string    `json:"emergency_name" gorm:"size:50;comment:紧急联系人姓名"`
	EmergencyTel   string    `json:"emergency_tel" gorm:"size:100;comment:紧急联系人电话"`
	Address        string    `json:"address" gorm:"size:200;comment:通讯地址"`
	Email          string    `json:"email" gorm:"size:50;comment:邮箱"`
	Source         int       `json:"source" gorm:"default:1;comment:注册来源：1系统录入2微信3web端4app5其它"`
	ValidDateBegin time.Time `json:"valid_date_begin" gorm:"comment:身份证有效期开始时间"`
	ValidDateEnd   time.Time `json:"valid_date_end" gorm:"comment:身份证有效期截止时间"`
	GraduatedFrom  string    `json:"graduated_from" gorm:"default:'';comment:毕业院校"`
	Schooling      string    `json:"schooling" gorm:"size:100;comment:学历"`
	DegreeNumber   string    `json:"degree_number" gorm:"size:100;comment:学位编号"`
	Professional   string    `json:"professional" gorm:"size:100;comment:职业"`
	Status         int       `json:"status" gorm:"not null;default:1;comment:用户行为状态：1正常，2已注销"`
	CreatedAt      int64     `json:"created_at" gorm:"autoCreateTime;comment:创建时间"`
	UpdatedAt      int64     `json:"updated_at" gorm:"autoCreateTime;comment:更新时间"`
	DeletedAt      *int64    `json:"deleted_at" gorm:"comment:删除时间"`
	Deleted        int       `json:"deleted" gorm:"not null;default:0;comment:是否删除"`
	Remark         string    `json:"remark" gorm:"default:'';comment:备注"`
}

var USER_PROFILE_NORMAL = base_model.UNIT_USER_PROFILE_NORMAL
var USER_PROFILE_CANCELED = base_model.UNIT_USER_PROFILE_CANCELED
var USER_CARD_TYPE_MAP = base_model.UNIT_CARD_TYPE_MAP

var USER_PROFILE_MAP = base_model.UNIT_USER_PROFILE_MAP

var USER_GENDER_MALE = base_model.UNIT_GENDER_MALE
var USER_GENDER_FEMALE = base_model.UNIT_GENDER_FEMALE
var USER_GENDER_MAP = base_model.UNIT_GENDER_MAP

var USER_CARD_TYPE_1 = base_model.UNIT_CARD_TYPE_1
var USER_CARD_TYPE_2 = base_model.UNIT_CARD_TYPE_2
var USER_CARD_TYPE_3 = base_model.UNIT_CARD_TYPE_3
var USER_CARD_TYPE_4 = base_model.UNIT_CARD_TYPE_4
var USER_CARD_TYPE_5 = base_model.UNIT_CARD_TYPE_5

var USER_SOURCE_SYSTEM = base_model.UNIT_USER_SOURCE_SYSTEM
var USER_SOURCE_WECHAT = base_model.UNIT_USER_SOURCE_WECHAT
var USER_SOURCE_WEB = base_model.UNIT_USER_SOURCE_WEB
var USER_SOURCE_APP = base_model.UNIT_USER_SOURCE_APP
var USER_SOURCE_OTHER = base_model.UNIT_USER_SOURCE_OTHER
var USER_SOURCE_MAP = base_model.UNIT_USER_SOURCE_MAP

func (m *UserProfile) TableName() string {
	return `user_profile`
}
