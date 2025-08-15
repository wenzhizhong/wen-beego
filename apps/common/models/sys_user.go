package models

import "time"

type SysUser struct {
	Id         string    `json:"id" gorm:"primaryKey;comment:用户ID"`
	Phone      string    `json:"phone" gorm:"unique;size:11;comment:手机号"`
	Password   string    `json:"password" gorm:"not null;comment:登录密码"`
	RoleType   int       `json:"role_type" gorm:"not null;default:4;comment:角色类型：1：管理员，2：产品测试，3：开发者，4：员工"`
	Name       string    `json:"name" gorm:"not null;size:20;comment:姓名"`
	EnName     string    `json:"en_name" gorm:"size:20;comment:英文名字"`
	Nickname   string    `json:"nickname" gorm:"size:20;comment:昵称"`
	Avatar     string    `json:"avatar" gorm:"size:256;comment:头像"`
	CreateTime time.Time `json:"create_time" gorm:"comment:创建时间"`
	UpdateTime time.Time `json:"update_time" gorm:"comment:更新时间"`
	Status     int       `json:"status" gorm:"not null;default:1;comment:状态：1：正常，2：禁用"`
	AuthKey    string    `json:"auth_key" gorm:"size:256;comment:验证key"`
}

func (m *SysUser) TableName() string {
	return "sys_user"
}
