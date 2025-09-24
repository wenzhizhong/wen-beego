package services

import (
	"WenBeego/apps/common/ar"
	"WenBeego/apps/common/global"
	"WenBeego/apps/common/helper"
	"WenBeego/apps/common/models"
	"errors"
	"strings"

	"WenBeego/apps/common/dto"

	"github.com/samber/lo"
)

type CommonAuth struct {
	userAr        ar.UserAr
	userProfileAr ar.UserProfileAr
}

// 登录
func (s *CommonAuth) Login(data dto.LoginDto, moduleName string) (loginInfo *dto.UserLoginInfoDto, err error) {
	err = s.checkLoginDto(&data)
	if err != nil {
		return
	}
	if err := s.checkAuthCode(data.AuthCode, data.AuthCodeId, data.AuthCodeType); err != nil {
		return nil, err
	}
	loginInfo, err2 := s.doLogin(data, moduleName)
	if err2 != nil {
		return nil, err2
	}

	return loginInfo, nil
}

// 注册
func (s *CommonAuth) Register(data dto.RegisterDto, moduleName string) (loginInfo *dto.UserLoginInfoDto, err error) {
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
	return loginInfo, nil
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

// 登录
func (s *CommonAuth) doLogin(data dto.LoginDto, moduleName string) (*dto.UserLoginInfoDto, error) {
	user, err := s.getAndCheckUser(data.Phone, data.Password)
	if err != nil {
		return nil, err
	}

	_, err = s.getAndCheckUserProfile(user.Id)
	if err != nil {
		return nil, err
	}

	loginInfo := &dto.UserLoginInfoDto{}
	if moduleName == "admin_plat" || moduleName == "admin_mchnt" {
		loginInfo, err = s.GetAdminLoginInfo(moduleName, user.Id)
	} else {
		loginInfo, err = s.GetApiLoginInfo(moduleName, user.Id)
	}

	if err != nil {
		return nil, err
	}
	return loginInfo, nil
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

// 获取并校验用户密码
func (s *CommonAuth) getAndCheckUser(phone string, password string) (*models.User, error) {
	user, err := s.userAr.GetByPhone(phone)
	if err != nil {
		if helper.DbNotFound(err) {
			return nil, errors.New("用户不存在")
		}
		global.Log.Error("获取用户异常: %v", err)
		return nil, errors.New("获取用户异常")
	}

	if err := helper.CheckPasswordRule(password); err != nil {
		return nil, err
	}

	if !helper.CompareCryptPassword(password, user.Password) {
		return nil, errors.New("账号或者密码错误")
	}
	return &user, nil
}

// 获取并校验用户信息
func (s *CommonAuth) getAndCheckUserProfile(userId string) (*models.UserProfile, error) {
	userProfile, err := s.userProfileAr.GetById(userId)
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

	return &userProfile, nil
}

/**
 * 获取后台管理用户信息
 */
func (s *CommonAuth) GetAdminLoginInfo(moduleName string, userId string) (*dto.UserLoginInfoDto, error) {
	defualtUnitId := ""
	perms := []string{}
	rolesClassifies := []string{}

	user, err := s.userAr.GetById(userId)
	if err != nil {
		return nil, err
	}
	userProfile, err := s.getAndCheckUserProfile(userId)
	if err != nil {
		return nil, err
	}
	defualtUnit, err := s.GetUserDefaultUnitId(moduleName, userId)
	if err != nil && !helper.DbNotFound(err) {
		return nil, err
	}
	if defualtUnit["id"] != nil {
		defualtUnitId = defualtUnit["id"].(string)
		rolesClassifies, _, err = s.getUserRolesClassifies(moduleName, defualtUnitId, userId)
		if err != nil {
			return nil, err
		}
		perms, err = s.GetUserPermissions(moduleName, defualtUnitId, userId)
		if err != nil {
			return nil, err
		}
	}
	exp, _ := global.GetConfigDiy("branca." + moduleName + ".exp")
	aud, _ := global.GetConfigDiy("branca." + moduleName + ".aud")
	cutTime := helper.GetTimestamp()
	brancaData := helper.BrancaData{}
	brancaData.Aud = aud.(string)
	brancaData.Iss = moduleName
	brancaData.Sub = user.Id
	brancaData.SubUnit = defualtUnitId
	brancaData.Role = strings.Join(rolesClassifies, ";")
	brancaData.Scope = strings.Join(perms, ";")
	brancaData.Exp = cutTime + int64(exp.(int))
	brancaData.Iat = cutTime
	token, err := helper.BrancaEncode(brancaData, moduleName)
	if err != nil {
		global.Log.Error("生成token异常: %v", err)
		return nil, errors.New("生成token异常")
	}
	refreshToken, err := helper.GetRefreshToken(moduleName, token, user.Id)
	if err != nil {
		global.Log.Error("生成refreshToken异常: %v", err)
		return nil, errors.New("生成refreshToken异常")
	}

	loginInfo := dto.UserLoginInfoDto{}
	loginInfo.UserInfo.User.Id = user.Id
	loginInfo.UserInfo.Name = user.Name
	loginInfo.UserInfo.Username = user.Name
	loginInfo.UserInfo.Phone = user.Phone[0:3] + "****" + user.Phone[7:11]
	loginInfo.UserInfo.Email = user.Email
	loginInfo.UserInfo.Avatar = userProfile.Avatar
	loginInfo.UserInfo.Expires = brancaData.Exp * 1000
	loginInfo.UserInfo.AccessToken = token
	loginInfo.UserInfo.RefreshToken = refreshToken
	loginInfo.UserInfo.DefaultUnitId = defualtUnitId
	loginInfo.UserInfo.Permissions = perms
	loginInfo.UserInfo.Roles = rolesClassifies
	loginInfo.UnitInfo = defualtUnit

	return &loginInfo, nil
}

/**
 * 获取普通用户登录信息
 */
func (s *CommonAuth) GetApiLoginInfo(moduleName string, userId string) (*dto.UserLoginInfoDto, error) {
	user, err := s.userAr.GetById(userId)
	if err != nil {
		return nil, err
	}
	userProfile, err := s.getAndCheckUserProfile(userId)
	if err != nil {
		return nil, err
	}

	exp, _ := global.GetConfigDiy("branca." + moduleName + ".exp")
	aud, _ := global.GetConfigDiy("branca." + moduleName + ".aud")
	cutTime := helper.GetTimestamp()
	brancaData := helper.BrancaData{}
	brancaData.Aud = aud.(string)
	brancaData.Iss = moduleName
	brancaData.Sub = user.Id
	brancaData.SubUnit = ""
	brancaData.Role = ""
	brancaData.Scope = ""
	brancaData.Exp = cutTime + int64(exp.(int))
	brancaData.Iat = cutTime
	token, err := helper.BrancaEncode(brancaData, moduleName)
	if err != nil {
		global.Log.Error("生成token异常: %v", err)
		return nil, errors.New("生成token异常")
	}
	refreshToken, err := helper.GetRefreshToken(moduleName, token, user.Id)
	if err != nil {
		global.Log.Error("生成refreshToken异常: %v", err)
		return nil, errors.New("生成refreshToken异常")
	}

	loginInfo := dto.UserLoginInfoDto{}
	loginInfo.UserInfo.User.Id = user.Id
	loginInfo.UserInfo.Name = user.Name
	loginInfo.UserInfo.Username = user.Name
	loginInfo.UserInfo.Phone = user.Phone[0:3] + "****" + user.Phone[7:11]
	loginInfo.UserInfo.Email = user.Email
	loginInfo.UserInfo.Avatar = userProfile.Avatar
	loginInfo.UserInfo.Expires = brancaData.Exp * 1000
	loginInfo.UserInfo.AccessToken = token
	loginInfo.UserInfo.RefreshToken = refreshToken
	loginInfo.UserInfo.DefaultUnitId = ""
	loginInfo.UserInfo.Permissions = []string{}
	loginInfo.UserInfo.Roles = []string{}

	return &loginInfo, nil
}

// 获取用户默认组织
func (s *CommonAuth) GetUserDefaultUnitId(moduleName string, userId string) (unitUserData map[string]interface{}, err error) {
	if moduleName == "admin_plat" {
		unitUserData, err = ar.GetUserDefaultUnit[*models.Plat, *models.PlatUser](userId)
	} else {
		unitUserData, err = ar.GetUserDefaultUnit[*models.Mchnt, *models.MchntUser](userId)
	}
	return unitUserData, err
}

/**
 * 获取用户角色分类
 * @param moduleName string 模块名称
 * @param unitId string 单位id
 * @param userId string 用户id
 * @return []string
 * @return error
 */
func (s *CommonAuth) getUserRolesClassifies(moduleName string, unitId string, userId string) (rolesClassifies []string, isAdmin bool, err error) {
	if unitId == "" {
		return
	}
	var dataList []map[string]interface{}
	if moduleName == "admin_plat" {
		dataList, err = ar.GetUserRoleClassifies(unitId, userId, &models.Plat{}, &models.PlatRole{}, &models.PlatRoleClassify{}, &models.PlatUserRole{})
	} else {
		dataList, err = ar.GetUserRoleClassifies(unitId, userId, &models.Mchnt{}, &models.MchntRole{}, &models.MchntRoleClassify{}, &models.MchntUserRole{})
	}

	for _, data := range dataList {
		roleName := data["name"].(string)
		rolesClassifies = append(rolesClassifies, roleName)
		if roleName == "admin" {
			isAdmin = true
			rolesClassifies = append(rolesClassifies[:0], "admin")
			break
		}
	}
	return
}

/**
 * 获取用户操作权限
 * @param moduleName 模块名称
 * @param unitId 组织ID
 * @param userId 用户ID
 * @return menuAuthList 用户权限列表
 *
 */
func (s *CommonAuth) GetUserPermissions(moduleName string, unitId string, userId string) (perms []string, err error) {

	var permissions []map[string]interface{}
	if moduleName == "admin_plat" {
		permissions, err = ar.GetUserPermissions(moduleName, unitId, userId, &models.PlatMenu{}, &models.PlatMenuPerms{}, &models.PlatRoleMenu{}, &models.PlatUserRole{}, &models.PlatRole{})
	} else {
		permissions, err = ar.GetUserPermissions(moduleName, unitId, userId, &models.MchntMenu{}, &models.MchntMenuPerms{}, &models.MchntRoleMenu{}, &models.MchntUserRole{}, &models.MchntRole{})
	}

	for _, permission := range permissions {
		perms = append(perms, permission["permission"].(string))
	}

	return perms, err
}

// 刷新token
func (s *CommonAuth) RefreshToken(moduleName string, brancaToken string, refreshToken string) (loginInfo *dto.UserLoginInfoDto, err error) {
	verifyRes, userId, _ := helper.VerifyRefreshToken(moduleName, brancaToken, refreshToken)
	if !verifyRes {
		return loginInfo, errors.New("refreshToken已过期，请重新登录")
	}

	if moduleName == "admin_plat" || moduleName == "admin_mchnt" {
		loginInfo, err = s.GetAdminLoginInfo(moduleName, userId)
	} else {
		loginInfo, err = s.GetApiLoginInfo(moduleName, userId)
	}

	if err != nil {
		return nil, err
	}
	return loginInfo, nil
}
