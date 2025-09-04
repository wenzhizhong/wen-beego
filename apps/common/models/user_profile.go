package models

import (
	"WenBeego/apps/common/global"
	"errors"
	"time"

	"gorm.io/gorm"
)

type UserProfile struct {
	Id                string         `json:"id" gorm:"not null;primaryKey;type:bpchar(36);comment:ID"`
	Avatar            string         `json:"avatar" gorm:"size:255;comment:头像"`
	CardType          int            `json:"card_type" gorm:"type:int2;comment:1大陆身份证2港澳台身份证3护照4军官证5其它"`
	CardNum           string         `json:"card_num" gorm:"size:100;comment:证件号码"`
	CardImages        string         `json:"card_images" gorm:"size:1000;comment:证件照片"`
	Sex               string         `json:"sex" gorm:"size:2;comment:性别:男，女"`
	BirthDate         time.Time      `json:"birth_date" gorm:"type:date;comment:出生日期"`
	Constellation     string         `json:"constellation" gorm:"size:50;comment:星座"`
	Occupation        string         `json:"occupation" gorm:"size:50;comment:职业"`
	Company           string         `json:"company" gorm:"size:500;comment:所属公司名称"`
	EmergencyName     string         `json:"emergency_name" gorm:"size:50;comment:紧急联系人姓名"`
	EmergencyTel      string         `json:"emergency_tel" gorm:"size:100;comment:紧急联系人电话"`
	Address           string         `json:"address" gorm:"size:200;comment:通讯地址"`
	EMail             string         `json:"e_mail" gorm:"size:50;comment:邮箱"`
	Source            string         `json:"source" gorm:"not null;default:'微信';comment:来源：微信,web,其它,app"`
	Headimgurl        string         `json:"headimgurl" gorm:"size:500;comment:头像"`
	ValidDateBegin    time.Time      `json:"valid_date_begin" gorm:"comment:身份证有效期开始时间"`
	ValidDateEnd      time.Time      `json:"valid_date_end" gorm:"comment:身份证有效期截止时间"`
	Schooling         string         `json:"schooling" gorm:"size:100;comment:学历"`
	DegreeNumber      string         `json:"degree_number" gorm:"size:100;comment:学位编号"`
	LearnProfessional string         `json:"learn_professional" gorm:"size:100;comment:所学专业"`
	Professional      string         `json:"professional" gorm:"size:100;comment:职业"`
	Status            int            `json:"status" gorm:"not null;default:1;comment:用户行为状态：1正常，2已注销"`
	CreatedAt         time.Time      `json:"created_at" gorm:"default:CURRENT_TIMESTAMP;comment:记录创建时间"`
	UpdatedAt         time.Time      `json:"updated_at" gorm:"default:CURRENT_TIMESTAMP;comment:记录修改时间"`
	DeletedAt         gorm.DeletedAt `json:"deleted_at" gorm:"comment:记录删除时间"`
	Deleted           bool           `json:"deleted" gorm:"not null;default:0;comment:是否删除"`
}

func (m *UserProfile) TableName() string {
	return `user_profile`
}

func (m *UserProfile) GetById(id string) (UserProfile, error) {
	user := UserProfile{}
	if id == "" {
		return user, errors.New("id不能为空")
	}
	result := global.GetReadDb().Where("id = ?", id).Take(&user)
	return user, result.Error
}
