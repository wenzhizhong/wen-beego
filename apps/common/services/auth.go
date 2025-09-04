package services

import (
	"WenBeego/apps/common/dto"
	"WenBeego/apps/common/global"
	"WenBeego/apps/common/helper"
	"WenBeego/apps/common/models"
	"errors"
	"strings"

	"github.com/samber/lo"
)

type CommonAuth struct {
	userModel   models.User
	userProfile models.UserProfile
}

// 登录
func (s *CommonAuth) Login(data dto.LoginDto, moduleName string) (userInfo *dto.UserInfoDto, err error) {
	err = s.checkLoginDto(&data)
	if err != nil {
		return
	}
	if err := s.checkAuthCode(data.AuthCode, data.AuthCodeId, data.AuthCodeType); err != nil {
		return nil, err
	}
	userInfo, err2 := s.doLogin(data, moduleName)
	if err2 != nil {
		return nil, err2
	}

	return userInfo, nil
}

// 注册
func (s *CommonAuth) Register(data dto.RegisterDto, moduleName string) (userInfo *dto.UserInfoDto, err error) {
	err = s.checkRegisterDto(data)
	if err != nil {
		return
	}
	if err := s.checkAuthCode(data.AuthCode, data.AuthCodeId, data.AuthCodeType); err != nil {
		return nil, err
	}

	// user := models.User{}
	// user.Phone = data.Phone
	// user.Name = data.Name
	// user.Email = data.Email
	// TODO: 注册
	return userInfo, nil
}

// 登录
func (s *CommonAuth) doLogin(data dto.LoginDto, moduleName string) (*dto.UserInfoDto, error) {
	user, err := s.userModel.GetByPhone(data.Phone)
	if err != nil {
		if helper.DbNotFound(err) {
			return nil, errors.New("用户不存在")
		}
		global.Log.Error("获取用户异常: %v", err)
		return nil, errors.New("获取用户异常")
	}

	if err := helper.CheckPasswordRule(data.Password); err != nil {
		return nil, err
	}

	if !helper.CompareCryptPassword(data.Password, user.Password) {
		return nil, errors.New("账号或者密码错误")
	}

	userProfile, err := s.userProfile.GetById(user.Id)
	if err != nil {
		if helper.DbNotFound(err) {
			return nil, errors.New("用户信息不存在")
		}
		global.Log.Error("获取用户信息异常: %v", err)
		return nil, errors.New("获取用户信息异常")
	}
	if userProfile.Status != 1 {
		return nil, errors.New("用户已注销")
	}

	brancaData := helper.BrancaData{}
	exp, _ := global.GetConfigDiy("branca.exp")
	brancaData.Aud = user.Id
	brancaData.Iss = moduleName
	brancaData.Sub = user.Id
	brancaData.Role = moduleName
	brancaData.Scope = moduleName
	brancaData.Exp = helper.GetTimestamp() + int64(exp.(int))
	token, _ := helper.BrancaEncode(brancaData)

	// userInfo := dto.UserInfoDto{User: user, UserProfile: userProfile}
	// userInfo.Password = ""
	// userInfo.CardImages = ""
	// userInfo.CardNum = ""
	userInfo := dto.UserInfoDto{}
	userInfo.User.Id = user.Id
	userInfo.Name = user.Name
	userInfo.Phone = user.Phone[0:3] + "****" + user.Phone[7:11]
	userInfo.Email = user.Email
	userInfo.Avatar = userProfile.Avatar
	userInfo.Token = token

	return &userInfo, nil
}

// 校验登录信息
func (s *CommonAuth) checkLoginDto(data *dto.LoginDto) error {
	data.Phone = strings.TrimSpace(data.Phone)
	data.Password = strings.TrimSpace(data.Password)

	if !helper.IsCellPhone(data.Phone) {
		return errors.New("手机号码格式错误")
	}

	if data.Password == "" {
		return errors.New("密码不能为空")
	}

	if data.AuthCode == "" || data.AuthCodeId == "" || data.AuthCodeType == "" {
		return errors.New("验证码不能为空")
	}

	return nil
}

// 校验注册信息
func (s *CommonAuth) checkRegisterDto(data dto.RegisterDto) error {
	data.Phone = strings.TrimSpace(data.Phone)
	data.Password = strings.TrimSpace(data.Password)

	if err := s.checkAuthCode(data.AuthCode, data.AuthCodeId, data.AuthCodeType); err != nil {
		return err
	}

	if !helper.IsCellPhone(data.Phone) {
		return errors.New("手机号码格式错误")
	}

	if err := helper.CheckPasswordRule(data.Password); err != nil {
		return err
	}
	return nil
}

// 校验验证码
func (s *CommonAuth) checkAuthCode(authCode string, authCodeId string, authCodeType string) error {
	authCode = strings.TrimSpace(authCode)
	authCodeType = strings.TrimSpace(authCodeType)

	if !lo.Contains(dto.AuthCodeTypes, authCodeType) {
		return errors.New("验证码类型错误")
	}

	switch authCodeType {
	case
		dto.AuthCodeTypeDigit,
		dto.AuthCodeTypeString,
		dto.AuthCodeTypeChinese,
		dto.AuthCodeTypeMath:

		if authCode == "" {
			return errors.New("验证码不能为空")
		}
		if !helper.VerifyCaptcha(authCodeType, authCodeId, authCode) {
			return errors.New("验证码错误")
		}
	case "sms":
		if authCode == "" {
			return errors.New("手机验证码不能为空")
		}
		// TODO: 短信验证
		return errors.New("短信验证码暂不支持")
	default:
		return errors.New("验证码类型错误")
	}

	return nil
}

// 获取验证码
func (s *CommonAuth) GetCatpcha(cpatchaType string) (data interface{}, err error) {
	id, b64s, _, err := helper.GetCaptcha(cpatchaType)
	data = map[string]interface{}{
		"id":   id,
		"b64s": b64s,
	}
	if err != nil {
		return
	}

	return
}
