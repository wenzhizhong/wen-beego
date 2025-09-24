package services

import (
	"WenBeego/apps/common/dto"
	commonServices "WenBeego/apps/common/services"
)

var cpatchaType = dto.AuthCodeTypeDigit

type Auth struct {
	commonAuth commonServices.CommonAuth
}

// 获取验证码
func (s *Auth) GetCatpcha() (interface{}, error) {
	return s.commonAuth.GetCatpcha(cpatchaType)
}

// 登录
func (s *Auth) Login(data dto.LoginDto, moduleName string) (interface{}, error) {
	data.AuthCodeType = cpatchaType
	return s.commonAuth.Login(data, moduleName)
}

// 注册
func (s *Auth) Register(data dto.RegisterDto, moduleName string) (interface{}, error) {
	data.AuthCodeType = cpatchaType
	return s.commonAuth.Register(data, moduleName)
}

// 刷新token
func (s *Auth) RefreshToken(moduleName string, brancaToken string, refreshToken string) (interface{}, error) {
	return s.commonAuth.RefreshToken(moduleName, brancaToken, refreshToken)
}
