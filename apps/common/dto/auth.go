package dto

var AuthCodeTypes = []string{"captcha", "email", "sms"}

type LoginDto struct {
	Phone        string `json:"phone"  validate:"required,numeric,len=11" example:"15912345678"`
	Password     string `json:"password" validate:"required,alphanum,min=6,max=20" example:"123456"`
	AuthCode     string `json:"authCode" example:"1234"`
	AuthCodeType string `json:"authCodeType" example:"captcha"`
}

type RegisterDto struct {
	Phone        string `json:"phone"  validate:"required,numeric,len=11" example:"15912345678"`
	Password     string `json:"password" validate:"required,alphanum,min=6,max=20" example:"123456"`
	AuthCode     string `json:"authCode" example:"1234"`
	AuthCodeType string `json:"authCodeType" example:"captcha"`
}
