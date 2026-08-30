package framework

import (
	"WenBeego/apps/common/helper"
	"WenBeego/apps/common/middleware/business_store"
	"WenBeego/apps/common/models"
	"WenBeego/apps/common/models/base_model"
	"WenBeego/apps/common/models_ar/base_ar"
	"errors"
	"strconv"
	"strings"
)

// 认证中间件服务层
type AuthMiddleware struct {
}

/**
 * 验证组织单位用户各种状态
 */
func (s *AuthMiddleware) CheckAuthAdminStatus(moduleName string, brancaData helper.BrancaData) (bool, error) {
	userId := brancaData.Sub
	unitId := brancaData.SubUnit
	unitUserId := brancaData.SubUnitUser
	if userId == "" || unitId == "" {
		return false, errors.New("CheckAuthAdminStatus(): 用户id,组织单位id，均不能为空。userId=" + userId + ", unitId=" + unitId)
	}

	type checkResult struct {
		status bool
		err    error
	}

	// 创建通道接收结果
	checkResultNum := 3
	ch := make(chan checkResult, checkResultNum)

	// 并发执行检查
	go func() {
		status, err := s.checkUnitUserProfileStatus(moduleName, userId, unitId)
		if !status && err == nil {
			err = errors.New("用户状态不可用")
		}
		ch <- checkResult{status, err}
	}()

	go func() {
		status, err := s.checkUnitStatus(moduleName, userId, unitId)
		if !status && err == nil {
			err = errors.New("组织单位状态不可用")
		}
		ch <- checkResult{status, err}
	}()

	go func() {
		status, err := s.checkUserRoleStatus(moduleName, unitUserId)
		if !status && err == nil {
			err = errors.New("角色状态不可用")
		}
		ch <- checkResult{status, err}
	}()

	// 收集结果并检查
	errorStr := ""
	for i := 0; i < checkResultNum; i++ {
		result := <-ch
		if result.err != nil {
			errorStr += result.err.Error() + ";\n"
		}
	}
	close(ch)
	if errorStr != "" {
		return false, errors.New(errorStr)
	}

	return true, nil
}

/**
 * 验证api用户各种状态
 */
func (s *AuthMiddleware) CheckAuthUserStatus(moduleName string, brancaData helper.BrancaData) (bool, error) {
	userId := brancaData.Sub
	if userId == "" {
		return false, errors.New("CheckAuthUserStatus(): 用户id，不能为空")
	}
	return s.checkUserProfileStatus(moduleName, userId)
}

func (s *AuthMiddleware) CheckAuthAdminRouters(moduleName string, brancaData helper.BrancaData, path string) (bool, error) {
	unitId := brancaData.SubUnit
	unitUserId := brancaData.SubUnitUser
	if unitUserId == "" || unitId == "" {
		return false, errors.New("CheckAuthAdminRouters(): 用户id,组织单位id，均不能为空。unitUserId=" + unitUserId + ", unitId=" + unitId)
	}

	isOk, err := s.checkUserRolePermissions(moduleName, unitUserId, unitId, path)
	return isOk, err
}

// 验证组织单位用户资料状态
func (s *AuthMiddleware) checkUnitUserProfileStatus(moduleName string, userId string, unitId string) (status bool, err error) {
	index := 0
	index, err = business_store.GetAumidUps(userId, unitId, 0)
	if err == nil && index > 0 {
		status = index == base_model.UNIT_USER_PROFILE_NORMAL
		if !status {
			err = errors.New("用户" + base_model.UNIT_USER_PROFILE_MAP[index])
		}
	}
	if err != nil {
		return
	}

	var data base_model.UnitUserProfile
	switch moduleName {
	case "admin_plat":
		data, err = base_ar.GetUserProfileOfUnitById[*models.PlatUser, *models.PlatUserProfile](userId, unitId)
	case "admin_mchnt":
		data, err = base_ar.GetUserProfileOfUnitById[*models.MchntUser, *models.MchntUserProfile](userId, unitId)
	default:
		err = errors.New("checkUnitUserProfileStatus:未知的模块名称")
	}
	if err != nil {
		return
	}
	if data.Id == "" {
		err = errors.New("用户资料信息不存在")
		return
	}

	if !(data.Status == base_model.UNIT_USER_PROFILE_NORMAL) {
		err = errors.New("用户" + base_model.UNIT_USER_PROFILE_MAP[data.Status])
		return
	}

	err = business_store.SetAumidUps(userId, unitId, data.Status)
	if err == nil {
		status = true
	}
	return
}

// 验证api用户资料状态
func (s *AuthMiddleware) checkUserProfileStatus(moduleName string, userId string) (status bool, err error) {
	index := 0
	index, err = business_store.GetAumidUps(userId, "", 0)
	if err == nil && index > 0 {
		status = index == models.USER_PROFILE_NORMAL
		if !status {
			err = errors.New("用户" + models.USER_PROFILE_MAP[index])
		}
	}
	if err != nil {
		return
	}

	var data models.UserProfile
	data, err = base_ar.GetUserProfileOfById(userId)
	if err != nil {
		return
	}
	if data.Id == "" {
		err = errors.New("用户资料信息不存在")
		return
	}

	if !(data.Status == models.USER_PROFILE_NORMAL) {
		err = errors.New("用户" + models.USER_PROFILE_MAP[data.Status])
		return
	}

	err = business_store.SetAumidUps(userId, "", data.Status)
	if err == nil {
		status = true
	}
	return
}

// 验证组织单位状态
func (s *AuthMiddleware) checkUnitStatus(moduleName string, userId string, unitId string) (status bool, err error) {
	index := -1
	index, err = business_store.GetAumidUs(userId, unitId, -1)
	if err == nil && index > -1 {
		switch moduleName {
		case "admin_plat":
			status = index == base_model.UNIT_STATUS_PASSED
			if !status {
				err = errors.New("用户" + base_model.UNIT_STATUS_MAP[index])
			}
		case "admin_mchnt":
			status = index == base_model.UNIT_STATUS_PASSED
			if !status {
				err = errors.New("用户" + base_model.UNIT_STATUS_MAP[index])
			}
		default:
			err = errors.New("checkUnitStatus:未知的模块名称")
		}
		return
	} else if err != nil {
		return
	}

	var data interface{}
	switch moduleName {
	case "admin_plat":
		data, err = base_ar.GetUserUnitById(userId, unitId, &models.Plat{}, &models.PlatUser{})
	case "admin_mchnt":
		data, err = base_ar.GetUserUnitById(userId, unitId, &models.Mchnt{}, &models.MchntUser{})
	default:
		err = errors.New("未知的模块名称")
	}
	if err != nil {
		return
	}
	if data == nil {
		err = errors.New("用户组织单位不存在")
		return
	}

	tmpStatus := -1
	switch v := data.(type) {
	case *models.Plat:
		tmpStatus = v.Status
		status = v.Status == models.USER_PROFILE_NORMAL
	case *models.Mchnt:
		tmpStatus = v.Status
		status = v.Status == models.USER_PROFILE_NORMAL
	default:
		err = errors.New("未知的用户单位类型")
	}
	if err != nil {
		return
	}
	if !status {
		err = errors.New("用户状态不可用, status=" + strconv.Itoa(tmpStatus))
		return
	}

	err = business_store.SetAumidUs(userId, unitId, tmpStatus)
	return
}

// 验证用户角色状态
func (s *AuthMiddleware) checkUserRoleStatus(moduleName string, unitUserId string) (status bool, err error) {

	index := -1
	index, err = business_store.GetAumidUrs(unitUserId, -1)
	if err == nil && index > -1 {
		status = index == base_model.UNIT_ROLE_STATUS_NORMAL
		if !status {
			err = errors.New("用户" + base_model.UNIT_ROLE_STATUS_MAP[index])
		}
		return
	}

	var roles []base_model.UnitRole
	switch moduleName {
	case "admin_plat":
		roles, err = base_ar.GetUserRole(moduleName, unitUserId, &models.PlatUserRole{}, &models.PlatRole{})
	case "admin_mchnt":
		roles, err = base_ar.GetUserRole(moduleName, unitUserId, &models.MchntUserRole{}, &models.MchntRole{})
	default:
		err = errors.New("未知的模块名称")
	}
	if err != nil {
		return
	}
	if len(roles) == 0 {
		err = errors.New("用户角色不存在")
		return
	}

	tmpStatus := -1
	for _, item := range roles {
		tmpStatus = item.Status
		status = item.Status == base_model.UNIT_ROLE_STATUS_NORMAL
		if status {
			break
		}
	}
	if !status {
		err = errors.New("用户角色状态不可用")
		return
	}

	err = business_store.SetAumidUrs(unitUserId, tmpStatus)
	return
}

func (s *AuthMiddleware) checkUserRolePermissions(moduleName string, unitUserId string, unitId string, path string) (status bool, err error) {
	tmpMap := map[string]bool{}
	exits, err := business_store.GetAumidUrp(moduleName, unitUserId, unitId)
	if err == nil && exits != "" {
		exitsArr := strings.Split(exits, ";")
		for _, item := range exitsArr {
			if item == "" {
				continue
			}
			tmpMap[item] = true
		}
	} else {
		var permissions []base_model.UnitMenu
		switch moduleName {
		case "admin_plat":
			permissions, err = base_ar.GetUserPermissions(moduleName, unitId, unitUserId, &models.PlatMenu{}, &models.PlatMenuMap{}, &models.PlatRoleMenu{}, &models.PlatUserRole{}, &models.PlatRole{})
		case "admin_mchnt":
			permissions, err = base_ar.GetUserPermissions(moduleName, unitId, unitUserId, &models.MchntMenu{}, &models.MchntMenuMap{}, &models.MchntRoleMenu{}, &models.MchntUserRole{}, &models.MchntRole{})
		default:
			err = errors.New("checkUserRolePermissions:未知的模块名称")
		}
		if err != nil {
			return
		}
		if len(permissions) == 0 {
			err = errors.New("用户权限不存在")
			return
		}

		keys := make([]string, 0)
		for _, item := range permissions {
			if item.Path == "" {
				continue
			}
			key := item.Path
			keys = append(keys, key)
			tmpMap[key] = true
		}

		err = business_store.SetAumidUrp(moduleName, unitUserId, unitId, strings.Join(keys, ";"))
		if err != nil {
			return
		}
	}
	if !tmpMap[path] {
		err = errors.New("没有权限：" + path)
		return
	}
	return true, nil
}
