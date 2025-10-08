package framework

import (
	"WenBeego/apps/common/helper"
	"WenBeego/apps/common/models"
	"WenBeego/apps/common/models/base_model"
	"WenBeego/apps/common/models_ar/base_ar"
	"errors"
	"strconv"
	"strings"
)

// 认证中间件服务层
type AuthMiddlewate struct {
}

/**
 * 验证用户各种状态
 */
func (s *AuthMiddlewate) CheckAuthAdminStatus(moduleName string, brancaData helper.BrancaData) (bool, error) {
	userId := brancaData.Sub
	unitId := brancaData.SubUnit
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
		status, err := s.checkUserRoleStatus(moduleName, userId, unitId)
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

func (s *AuthMiddlewate) CheckAuthAdminRouters(moduleName string, brancaData helper.BrancaData, path string) (bool, error) {
	userId := brancaData.Sub
	unitId := brancaData.SubUnit
	if userId == "" || unitId == "" {
		return false, errors.New("CheckAuthAdminRouters(): 用户id,组织单位id，均不能为空。userId=" + userId + ", unitId=" + unitId)
	}

	isOk, err := s.checkUserRolePermissions(moduleName, userId, unitId, path)
	return isOk, err
}

// 验证用户资料状态
func (s *AuthMiddlewate) checkUnitUserProfileStatus(moduleName string, userId string, unitId string) (status bool, err error) {
	redisKey := "AUMID_UPS:" + helper.Md5(userId+unitId)

	exits, err := helper.RedisGet(redisKey)
	if err == nil && exits != "" {
		index, _ := strconv.Atoi(exits)
		status = exits == strconv.Itoa(base_model.UNIT_USER_PROFILE_NORMAL)
		if !status {
			err = errors.New("用户" + base_model.UNIT_USER_PROFILE_MAP[index])
		}
		return
	} else if err != nil {
		return
	}

	var data base_model.UnitUserProfile
	if moduleName == "admin_plat" {
		data, err = base_ar.GetUserProfileOfUnitById[*models.PlatUser, *models.PlatUserProfile](userId, unitId)
	} else if moduleName == "admin_mchnt" {
		data, err = base_ar.GetUserProfileOfUnitById[*models.MchntUser, *models.MchntUserProfile](userId, unitId)
	} else {
		err = errors.New("未知的模块名称")
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

	err = helper.RedisPut(redisKey, data.Status, 4*60*60)
	return
}

// 验证组织单位状态
func (s *AuthMiddlewate) checkUnitStatus(moduleName string, userId string, unitId string) (status bool, err error) {
	redisKey := "AUMID_US:" + helper.Md5(userId+unitId)
	exits, err := helper.RedisGet(redisKey)
	if err == nil && exits != "" {
		index, _ := strconv.Atoi(exits)
		if moduleName == "admin_plat" {
			status = exits == strconv.Itoa(base_model.UNIT_STATUS_PASSED)
			if !status {
				err = errors.New("用户" + base_model.UNIT_STATUS_MAP[index])
			}
		} else if moduleName == "admin_mchnt" {
			status = exits == strconv.Itoa(base_model.UNIT_STATUS_PASSED)
			if !status {
				err = errors.New("用户" + base_model.UNIT_STATUS_MAP[index])
			}
		} else {
			err = errors.New("未知的模块名称")
		}
		return
	} else if err != nil {
		return
	}

	var data interface{}
	if moduleName == "admin_plat" {
		data, err = base_ar.GetUserUnitById(userId, unitId, &models.Plat{}, &models.PlatUser{})
	} else if moduleName == "admin_mchnt" {
		data, err = base_ar.GetUserUnitById(userId, unitId, &models.Mchnt{}, &models.MchntUser{})
	} else {
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
		status = v.Status == models.USER_STATUS_NORMAL
	case *models.Mchnt:
		tmpStatus = v.Status
		status = v.Status == models.USER_STATUS_NORMAL
	default:
		err = errors.New("未知的用户单位类型")
	}
	if err != nil {
		return
	}
	if !status {
		err = errors.New("用户状态不可用")
		return
	}
	err = helper.RedisPut(redisKey, strconv.Itoa(tmpStatus), 4*60*60)
	return
}

// 验证用户角色状态
func (s *AuthMiddlewate) checkUserRoleStatus(moduleName string, userId string, unitId string) (status bool, err error) {
	redisKey := "AUMID_URS:" + helper.Md5(userId+unitId)
	exits, err := helper.RedisGet(redisKey)
	if err == nil && exits != "" {
		index, _ := strconv.Atoi(exits)

		status = exits == strconv.Itoa(base_model.UNIT_ROLE_STATUS_NORMAL)
		if !status {
			err = errors.New("用户" + base_model.UNIT_ROLE_STATUS_MAP[index])
		}
		return
	}

	var roles []base_model.UnitRole
	if moduleName == "admin_plat" {
		roles, err = base_ar.GetUserRole(moduleName, unitId, userId, &models.PlatUserRole{}, &models.PlatRole{})
	} else if moduleName == "admin_mchnt" {
		roles, err = base_ar.GetUserRole(moduleName, unitId, userId, &models.MchntUserRole{}, &models.MchntRole{})
	} else {
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

	err = helper.RedisPut(redisKey, strconv.Itoa(tmpStatus), 4*60*60)
	return
}

func (s *AuthMiddlewate) checkUserRolePermissions(moduleName string, userId string, unitId string, path string) (status bool, err error) {
	redisKey := "AUMID_URP:" + helper.Md5(moduleName+userId+unitId)

	tmpMap := map[string]bool{}
	exits, err := helper.RedisGet(redisKey)
	if err == nil && exits != "" {
		exitsArr := strings.Split(exits, ";")
		for _, item := range exitsArr {
			tmpMap[item] = true
		}
	} else {
		var permissions []base_model.UnitMenuPerms
		if moduleName == "admin_plat" {
			permissions, err = base_ar.GetUserPermissions(moduleName, unitId, userId, &models.PlatMenu{}, &models.PlatMenuPerms{}, &models.PlatRoleMenu{}, &models.PlatUserRole{}, &models.PlatRole{})
		} else if moduleName == "admin_mchnt" {
			permissions, err = base_ar.GetUserPermissions(moduleName, unitId, userId, &models.MchntMenu{}, &models.MchntMenuPerms{}, &models.MchntRoleMenu{}, &models.MchntUserRole{}, &models.MchntRole{})
		} else {
			err = errors.New("未知的模块名称")
		}
		if err != nil {
			return
		}
		if len(permissions) == 0 {
			err = errors.New("用户权限不存在")
			return
		}

		keys := make([]string, len(permissions))
		for _, item := range permissions {
			key := item.Uri
			keys = append(keys, key)
			tmpMap[key] = true
		}

		err = helper.RedisPut(redisKey, strings.Join(keys, ";"), 4*60*60)
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
