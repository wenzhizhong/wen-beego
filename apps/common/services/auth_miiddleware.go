package services

import (
	"WenBeego/apps/common/ar"
	"WenBeego/apps/common/helper"
	"WenBeego/apps/common/models"
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
		if moduleName == "admin_plat" {
			status = exits == strconv.Itoa(models.PLAT_USER_PROFILE_STATUS_NORMAL)
			if !status {
				err = errors.New("用户" + models.PLAT_USER_PROFILE_STATUS_MAP[index])
			}
		} else {
			status = exits == strconv.Itoa(models.MCHNT_USER_PROFILE_STATUS_NORMAL)
			if !status {
				err = errors.New("用户" + models.MCHNT_USER_PROFILE_STATUS_MAP[index])
			}
		}
		return
	} else if err != nil {
		return
	}

	var data interface{}
	if moduleName == "admin_plat" {
		data, err = ar.GetUserProfileOfUnitById[*models.PlatUser, *models.PlatUserProfile](userId, unitId)
	} else {
		data, err = ar.GetUserProfileOfUnitById[*models.MchntUser, *models.MchntUserProfile](userId, unitId)
	}
	if err != nil {
		return
	}
	if data == nil {
		err = errors.New("用户资料信息不存在")
		return
	}

	tmpStatus := -1
	switch v := data.(type) {
	case *models.PlatUserProfile:
		tmpStatus = v.Status
		status = v.Status == models.PLAT_USER_PROFILE_STATUS_NORMAL
		if !status {
			err = errors.New("用户" + models.PLAT_USER_PROFILE_STATUS_MAP[v.Status])
		}
	case *models.MchntUserProfile:
		tmpStatus = v.Status
		status = v.Status == models.MCHNT_USER_PROFILE_STATUS_NORMAL
		if !status {
			err = errors.New("用户" + models.MCHNT_USER_PROFILE_STATUS_MAP[v.Status])
		}
	default:
		err = errors.New("未知的单位用户信息类型")
	}
	if err != nil {
		return
	}

	err = helper.RedisPut(redisKey, tmpStatus, 4*60*60)
	return
}

// 验证组织单位状态
func (s *AuthMiddlewate) checkUnitStatus(moduleName string, userId string, unitId string) (status bool, err error) {
	redisKey := "AUMID_US:" + helper.Md5(userId+unitId)
	exits, err := helper.RedisGet(redisKey)
	if err == nil && exits != "" {
		index, _ := strconv.Atoi(exits)
		if moduleName == "admin_plat" {
			status = exits == strconv.Itoa(models.PLAT_STATUS_PASSED)
			if !status {
				err = errors.New("用户" + models.PLAT_STATUS_MAP[index])
			}
		} else {
			status = exits == strconv.Itoa(models.MCHNT_STATUS_PASSED)
			if !status {
				err = errors.New("用户" + models.MCHNT_STATUS_MAP[index])
			}
		}
		return
	} else if err != nil {
		return
	}

	var data interface{}
	if moduleName == "admin_plat" {
		data, err = ar.GetUserUnitById(userId, unitId, &models.Plat{}, &models.PlatUser{})
	} else {
		data, err = ar.GetUserUnitById(userId, unitId, &models.Mchnt{}, &models.MchntUser{})
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
		if moduleName == "admin_plat" {
			status = exits == strconv.Itoa(models.PLAT_ROLE_STATUS_NORMAL)
			if !status {
				err = errors.New("用户" + models.PLAT_ROLE_STATUS_MAP[index])
			}
		} else {
			status = exits == strconv.Itoa(models.MCHNT_ROLE_STATUS_NORMAL)
			if !status {
				err = errors.New("用户" + models.MCHNT_ROLE_STATUS_MAP[index])
			}
		}
		return
	}

	var datas []interface{}
	if moduleName == "admin_plat" {
		result, tmpErr := ar.GetUserRole(moduleName, unitId, userId, &models.PlatUserRole{}, &models.PlatRole{})
		err = tmpErr
		if err == nil && result != nil {
			for _, item := range result {
				datas = append(datas, item)
			}
		}

	} else {
		result, tmpErr := ar.GetUserRole(moduleName, unitId, userId, &models.MchntUserRole{}, &models.MchntRole{})
		err = tmpErr
		if err == nil && result != nil {
			for _, item := range result {
				datas = append(datas, item)
			}
		}
	}
	if err != nil {
		return
	}
	if len(datas) == 0 {
		err = errors.New("用户角色不存在")
		return
	}

	tmpStatus := -1
	for _, item := range datas {
		switch v := item.(type) {
		case *models.PlatRole:
			tmpStatus = v.Status
			status = v.Status == models.PLAT_ROLE_STATUS_NORMAL
		case *models.MchntRole:
			tmpStatus = v.Status
			status = v.Status == models.MCHNT_ROLE_STATUS_NORMAL
		}
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
		var permissions []map[string]interface{}
		if moduleName == "admin_plat" {
			permissions, err = ar.GetUserPermissions(moduleName, unitId, userId, &models.PlatMenu{}, &models.PlatMenuPerms{}, &models.PlatRoleMenu{}, &models.PlatUserRole{}, &models.PlatRole{})
		} else {
			permissions, err = ar.GetUserPermissions(moduleName, unitId, userId, &models.MchntMenu{}, &models.MchntMenuPerms{}, &models.MchntRoleMenu{}, &models.MchntUserRole{}, &models.MchntRole{})
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
			key := item["uri"].(string)
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
