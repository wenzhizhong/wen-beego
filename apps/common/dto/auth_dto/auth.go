package auth_dto

import "WenBeego/apps/common/models"

var AuthCodeTypeDigit = "captcha-digit"     // "数字验证码"
var AuthCodeTypeString = "captcha-string"   // "字符验证码"
var AuthCodeTypeMath = "captcha-math"       // "数学验证码"
var AuthCodeTypeChinese = "captcha-chinese" // "中文验证码"
var AuthCodeTypeEmail = "email"             // "邮箱验证码"
var AuthCodeTypeSms = "sms"                 // "短信验证码"

var AuthCodeTypes = []string{"captcha-digit", "captcha-string", "captcha-math", "captcha-chinese", "email", "sms"}

// 登录
type LoginDto struct {
	Phone        string `json:"phone"  validate:"required,numeric,len=11" example:"15912345678"`
	Password     string `json:"password" validate:"required,alphanum,min=6,max=20" example:"123456"`
	AuthCode     string `json:"authCode" example:"1234"`
	AuthCodeId   string `json:"authCodeId" example:"J4HGcl1gCLIdPm6T"`
	AuthCodeType string `json:"authCodeType" example:"captcha"`
}

// 注册
type RegisterDto struct {
	Phone        string `json:"phone"  validate:"required,numeric,len=11" example:"15912345678"`
	Password     string `json:"password" validate:"required,alphanum,min=6,max=20" example:"123456"`
	AuthCode     string `json:"authCode" example:"1234"`
	AuthCodeId   string `json:"authCodeId" example:"J4HGcl1gCLIdPm6T"`
	AuthCodeType string `json:"authCodeType" example:"captcha"`
}

// 用户信息
type UserLoginInfoDto struct {
	UserInfo struct {
		models.User
		models.UserProfile
		Expires       int64    `json:"expires"`
		AccessToken   string   `json:"accessToken"`
		RefreshToken  string   `json:"refreshToken"`
		DefaultUnitId string   `json:"default_unit_id"`
		Roles         []string `json:"roles"`
		Permissions   []string `json:"permissions"`
	} `json:"userInfo"`
	UnitInfo interface{} `json:"unitInfo"`
}

// 刷新token
type RefreshTokenDto struct {
	RefreshToken string `json:"refreshToken"`
	BrancaToken  string `json:"accessToken"`
}
