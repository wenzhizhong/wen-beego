package auth

import (
	"WenBeego/apps/common/dto/auth_dto"
	commonServices "WenBeego/apps/common/services/auth"
)

var cpatchaType = auth_dto.AuthCodeTypeDigit

type Auth struct {
	commonAuth commonServices.CommonAuth
}

// 获取验证码
func (s *Auth) GetCatpcha() (interface{}, error) {
	return s.commonAuth.GetCatpcha(cpatchaType)
}

// 登录
func (s *Auth) Login(data auth_dto.LoginDto, moduleName, host string) (interface{}, error) {
	data.AuthCodeType = cpatchaType
	return s.commonAuth.Login(data, moduleName, host)
}

// 注册
func (s *Auth) Register(data auth_dto.RegisterDto, moduleName string) (interface{}, error) {
	data.AuthCodeType = cpatchaType
	return s.commonAuth.Register(data, moduleName)
}

// 刷新token
func (s *Auth) RefreshToken(moduleName string, brancaToken string, refreshToken string) (interface{}, error) {
	return s.commonAuth.RefreshToken(moduleName, brancaToken, refreshToken)
}
