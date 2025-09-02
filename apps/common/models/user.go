package models

import (
	"WenBeego/apps/common/global"
	"errors"
)

type User struct {
	Id           string `json:"id" gorm:"not null;primaryKey;comment:用户ID"`
	DefaultOrgId string `json:"default_org_id" gorm:"not null; comment:默认组织ID"`
	Phone        string `json:"phone" gorm:"not null;unique;size:11;comment:手机号"`
	Name         string `json:"name" gorm:"not null;size:20;comment:姓名"`
	Email        string `json:"email" gorm:"size:64;comment:验证key"`
	Password     string `json:"password" gorm:"not null;comment:登录密码"`
	WxOpenid     string `json:"wx_openid" gorm:"size:64;comment:微信openid"`
}

func (m *User) TableName() string {
	return "user"
}

func (m *User) GetById(id string) (User, error) {
	user := User{}
	if id == "" {
		return user, errors.New("id不能为空")
	}
	result := global.GetReadDb().Where("id = ?", id).Take(&user)
	return user, result.Error
}

func (m *User) GetByPhone(phone string) (User, error) {
	user := User{}
	if phone == "" {
		return user, errors.New("手机号不能为空")
	}
	result := global.GetReadDb().Where("phone = ?", phone).Take(&user)
	return user, result.Error
}
