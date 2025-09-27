package models

type User struct {
	Id       string `json:"id" gorm:"not null;primaryKey;comment:用户ID"`
	Phone    string `json:"phone" gorm:"not null;unique;size:11;comment:手机号"`
	Name     string `json:"name" gorm:"not null;size:20;comment:姓名"`
	Username string `json:"username" gorm:"-"`
	Email    string `json:"email" gorm:"size:64;comment:验证key"`
	Password string `json:"password" gorm:"not null;comment:登录密码"`
	// WxOpenid string `json:"wx_openid" gorm:"size:64;comment:微信openid"`
}

func (m *User) TableName() string {
	return `user`
}
