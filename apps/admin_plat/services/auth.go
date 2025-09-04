package services

import (
	"WenBeego/apps/common/dto"
	commonServices "WenBeego/apps/common/services"
)

var cpatchaType = dto.AuthCodeTypeDigit

type Auth struct {
	CommonAuth commonServices.CommonAuth
}

// 获取验证码
func (s *Auth) GetCatpcha() (interface{}, error) {
	return s.CommonAuth.GetCatpcha(cpatchaType)
}

// 登录
func (s *Auth) Login(data dto.LoginDto, moduleName string) (interface{}, error) {
	data.AuthCodeType = cpatchaType
	return s.CommonAuth.Login(data, moduleName)
}

// 注册
func (s *Auth) Register(data dto.RegisterDto, moduleName string) (interface{}, error) {
	data.AuthCodeType = cpatchaType
	return s.CommonAuth.Register(data, moduleName)
}
