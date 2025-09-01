package services

import (
	"WenBeego/apps/common/dto"
	"WenBeego/apps/common/helper"
	"errors"
	"strings"

	"github.com/samber/lo"
)

type Auth struct {
}

func (s *Auth) Login(data dto.LoginDto) (interface{}, error) {
	err := s.checkLoginDto(&data)
	if err != nil {
		return "", err
	}
	// TODO
	return "123456", nil
}
func (s *Auth) Register(data dto.RegisterDto) (interface{}, error) {
	err := s.checkRegisterDto(data)
	if err != nil {
		return "", err
	}
	// TODO
	return "123456", nil
}

func (s *Auth) checkLoginDto(data *dto.LoginDto) error {
	data.Phone = strings.TrimSpace(data.Phone)
	data.Password = strings.TrimSpace(data.Password)

	if !helper.IsCellPhone(data.Phone) {
		return errors.New("手机号码格式错误")
	}

	if data.Password == "" {
		return errors.New("密码不能为空")
	}

	return nil
}

func (s *Auth) checkRegisterDto(data dto.RegisterDto) error {
	data.Phone = strings.TrimSpace(data.Phone)
	data.Password = strings.TrimSpace(data.Password)

	if err := s.checkAuthCode(data); err != nil {
		return err
	}

	if !helper.IsCellPhone(data.Phone) {
		return errors.New("手机号码格式错误")
	}

	if err := helper.ValidatePassword(data.Password); err != nil {
		return err
	}
	return nil
}

func (s *Auth) checkAuthCode(data dto.RegisterDto) error {
	data.AuthCode = strings.TrimSpace(data.AuthCode)
	data.AuthCodeType = strings.TrimSpace(data.AuthCodeType)

	if !lo.Contains(dto.AuthCodeTypes, data.AuthCodeType) {
		return errors.New("验证码类型错误")
	}

	switch data.AuthCodeType {
	case "captcha":
		if data.AuthCode == "" {
			return errors.New("验证码不能为空")
		}
		// TODO: 验证码验证
	case "sms":
		if data.AuthCode == "" {
			return errors.New("手机验证码不能为空")
		}
		// TODO: 验证码验证
	default:
		return errors.New("验证码类型错误")
	}

	return nil
}
