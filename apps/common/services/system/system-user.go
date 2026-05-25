package system

import (
	dto "WenBeego/apps/common/dto_vo"
	"WenBeego/apps/common/dto_vo/page_dto"
	"WenBeego/apps/common/dto_vo/unit_dto"
	"WenBeego/apps/common/dto_vo/user_dto"
	"WenBeego/apps/common/dto_vo/user_vo"
	"WenBeego/apps/common/global"
	"WenBeego/apps/common/helper"
	"WenBeego/apps/common/middleware/business_store"
	"WenBeego/apps/common/models"
	"WenBeego/apps/common/models/base_model"
	"WenBeego/apps/common/models/itf"
	"WenBeego/apps/common/models_ar/base_ar"
	"errors"
	"fmt"
	"strings"

	"gorm.io/gorm"
)

type User struct {
}

// 系统管理-获取用户列表
func (s *User) GetUserList(reqDto page_dto.SystemUserListReqDto) (resultDto dto.RespDataListDto, err error) {
	data := make([]user_vo.SystemUserListVo, 0)
	var count int64 = 0

	if res, err1 := helper.CheckUserHasUnit(reqDto.ModuleName, reqDto.UserId, reqDto.SelectUnitIds); !res {
		err = helper.Ternary(err1 != nil, err1, errors.New("GetUserList：用户没有组织单位权限"))
		return
	}

	switch reqDto.ModuleName {
	case "admin_plat":
		data, count, err = base_ar.GetUserListOfUnitById(reqDto, &models.PlatUser{}, &models.PlatUserProfile{}, &models.PlatDept{}, &models.PlattUserDept{}, &models.PlatRole{}, &models.PlatUserRole{})
	case "admin_mchnt":
		data, count, err = base_ar.GetUserListOfUnitById(reqDto, &models.MchntUser{}, &models.MchntUserProfile{}, &models.MchntDept{}, &models.MchntUserDept{}, &models.MchntRole{}, &models.MchntUserRole{})
	default:
		err = errors.New("GetUserList：模块名称错误")
	}
	if err != nil {
		return
	}

	for k, v := range data {
		tmpId := helper.Ternary(data[k].Id != "", data[k].Id, data[k].UnitUser.Id)
		if tmpId != "" && data[k].UnitUserProfile.Id == "" {
			data[k].UnitUserProfile.Id = tmpId
		}
		data[k].AvatarLink, _ = helper.LocalFileSign(reqDto.Host, v.Avatar)
	}
	resultDto, err = helper.GetRespDataListDto(reqDto.PageSize, reqDto.CurrentPage, count, data)
	return
}

func (s *User) GetUserListForAdminPlat(reqDto page_dto.SystemUserListReqDto) (resultDto dto.RespDataListDto, err error) {
	data := make([]user_vo.SystemUserListVo, 0)
	var count int64 = 0

	// if res, err1 := helper.CheckUserHasUnit(reqDto.ModuleName, reqDto.UserId, reqDto.SelectUnitIds); !res {
	// 	err = helper.Ternary(err1 != nil, err1, errors.New("GetUserList：用户没有组织单位权限"))
	// 	return
	// }

	data, count, err = base_ar.GetUserListOfUnitById(reqDto, &models.MchntUser{}, &models.MchntUserProfile{}, &models.MchntDept{}, &models.MchntUserDept{}, &models.MchntRole{}, &models.MchntUserRole{})
	if err != nil {
		return
	}

	for k, v := range data {
		tmpId := helper.Ternary(data[k].Id != "", data[k].Id, data[k].UnitUser.Id)
		if tmpId != "" && data[k].UnitUserProfile.Id == "" {
			data[k].UnitUserProfile.Id = tmpId
		}
		data[k].AvatarLink, _ = helper.LocalFileSign(reqDto.Host, v.Avatar)
	}
	resultDto, err = helper.GetRespDataListDto(reqDto.PageSize, reqDto.CurrentPage, count, data)
	return
}

// 获取可选角色树
func (s *User) GetUnitRoleTree(baseParamDto dto.BaseParamDto, selectUnitIds []string) (data interface{}, err error) {
	dataList := make([]base_model.UnitRole, 0)

	if res, err1 := helper.CheckUserHasUnit(baseParamDto.ModuleName, baseParamDto.UserId, selectUnitIds); !res {
		err = helper.Ternary(err1 != nil, err1, errors.New("GetUserList：用户没有组织单位权限"))
		return
	}

	switch baseParamDto.ModuleName {
	case "admin_plat":
		dataList, err = base_ar.GetUnitRoleTree(selectUnitIds, &models.PlatRole{})
	case "admin_mchnt":
		dataList, err = base_ar.GetUnitRoleTree(selectUnitIds, &models.MchntRole{})
	default:
		err = errors.New("GetUnitRoleTree:模块名称错误")
	}
	if err != nil {
		return
	}

	data = struct {
		List interface{} `json:"list"`
	}{
		List: dataList,
	}
	return
}
func (s *User) GetUnitRoleTreeForAdminPlat(baseParamDto dto.BaseParamDto, selectUnitIds []string) (data interface{}, err error) {
	dataList := make([]base_model.UnitRole, 0)

	// if res, err1 := helper.CheckUserHasUnit(baseParamDto.ModuleName, baseParamDto.UserId, selectUnitIds); !res {
	// 	err = helper.Ternary(err1 != nil, err1, errors.New("GetUserList：用户没有组织单位权限"))
	// 	return
	// }

	dataList, err = base_ar.GetUnitRoleTree(selectUnitIds, &models.MchntRole{})
	if err != nil {
		return
	}

	data = struct {
		List interface{} `json:"list"`
	}{
		List: dataList,
	}
	return
}

// 保存用户
func (s *User) SaveUser(baseParamDto dto.BaseParamDto, unitUserSaveDto *user_dto.UnitUserSaveDto) (data interface{}, err error) {
	var unitUserModel itf.UnitUserItf
	var unitUserProfileModel itf.UserProfileItf
	var unitDeptModel itf.DeptItf
	var unitRoleModel itf.RoleItf
	var unitUserDeptModel itf.UserDeptItf
	var unitUserRoleModel itf.UserRoleItf

	switch baseParamDto.ModuleName {
	case "admin_plat":
		unitUserModel = &models.PlatUser{}
		unitUserProfileModel = &models.PlatUserProfile{}
		unitDeptModel = &models.PlatDept{}
		unitRoleModel = &models.PlatRole{}
		unitUserDeptModel = &models.PlattUserDept{}
		unitUserRoleModel = &models.PlatUserRole{}
	case "admin_mchnt":
		unitUserModel = &models.MchntUser{}
		unitUserProfileModel = &models.MchntUserProfile{}
		unitDeptModel = &models.MchntDept{}
		unitRoleModel = &models.MchntRole{}
		unitUserDeptModel = &models.MchntUserDept{}
		unitUserRoleModel = &models.MchntUserRole{}
	default:
		err = errors.New("SaveUser:模块名称错误")
	}
	if err != nil {
		return
	}

	isAddUnitUser := unitUserSaveDto.UnitUserDto.Id == ""
	// 检查请求数据
	err = s.checkRequestData(baseParamDto, unitUserSaveDto)
	if err != nil {
		return
	}
	// 检查组织结构
	err = s.CheckOrgStructure(baseParamDto.UnitUserId, unitUserSaveDto.UnitUserDto.UnitId, unitUserSaveDto.DeptId, unitUserSaveDto.RoleId, unitUserModel, unitDeptModel, unitRoleModel)
	if err != nil {
		return
	}
	// 检查用户数据
	userData, unitUserData, err1 := s.checkAndSetUserData(unitUserSaveDto, unitUserModel, unitUserProfileModel)
	if err1 != nil {
		err = err1
		return
	}

	unitUserId := ""
	err = global.WriteDb.Transaction(func(tx *gorm.DB) (err error) {
		unitUserId, err = s.doSaveUser(tx, isAddUnitUser, baseParamDto, unitUserSaveDto, userData, unitUserData, unitUserModel, unitUserProfileModel, unitDeptModel, unitUserDeptModel, unitRoleModel, unitUserRoleModel)
		return
	})

	// 清空用户权限认证缓存
	business_store.ClearAumid()

	data = struct {
		UnitUserId string `json:"unitUserId"`
	}{
		UnitUserId: unitUserId,
	}
	return
}
func (s *User) SaveUserForAdminPlat(baseParamDto dto.BaseParamDto, unitUserSaveDto *user_dto.UnitUserSaveDto) (data interface{}, err error) {
	if true {
		return data, errors.New("平台没有操作权限")
	}
	return
}

func (s *User) doSaveUser(
	tx *gorm.DB,
	isAddUnitUser bool,
	baseParamDto dto.BaseParamDto,
	unitUserSaveDto *user_dto.UnitUserSaveDto,
	userData user_dto.UserAllDataDto,
	unitUserData unit_dto.UnitUserAllDataDto,
	unitUserModel itf.UnitUserItf,
	unitUserProfileModel itf.UserProfileItf,
	unitDeptModel itf.DeptItf,
	unitUserDeptModel itf.UserDeptItf,
	unitRoleModel itf.RoleItf,
	unitUserRoleModel itf.UserRoleItf,
) (unitUserId string, err error) {
	if isAddUnitUser {
		unitUserId, err = s.doInsertUser(tx, baseParamDto, unitUserSaveDto, userData, unitUserData, unitUserModel, unitUserProfileModel, unitDeptModel, unitUserDeptModel, unitRoleModel, unitUserRoleModel)
	} else {
		unitUserId, err = s.doEditUser(tx, baseParamDto, unitUserSaveDto, userData, unitUserData, unitUserModel, unitUserProfileModel, unitDeptModel, unitUserDeptModel, unitRoleModel, unitUserRoleModel)
	}
	if err != nil {
		return
	}

	err = base_ar.SaveUnitUserRole(tx, unitUserSaveDto.UnitUserDto.UnitId, unitUserSaveDto.UnitUserRoleDto, unitRoleModel, unitUserRoleModel)
	if err != nil {
		return
	}
	err = base_ar.SaveUnitUserDept(tx, unitUserSaveDto.UnitUserDto.UnitId, unitUserSaveDto.UnitUserDeptDto, unitDeptModel, unitUserDeptModel)
	if err != nil {
		return
	}

	return
}

// 系统录入-新增用户
func (s *User) doInsertUser(
	tx *gorm.DB,
	baseParamDto dto.BaseParamDto,
	unitUserSaveDto *user_dto.UnitUserSaveDto,
	userData user_dto.UserAllDataDto,
	unitUserData unit_dto.UnitUserAllDataDto,
	unitUserModel itf.UnitUserItf,
	unitUserProfileModel itf.UserProfileItf,
	unitDeptModel itf.DeptItf,
	unitUserDeptModel itf.UserDeptItf,
	unitRoleModel itf.RoleItf,
	unitUserRoleModel itf.UserRoleItf,
) (unitUserId string, err error) {
	// user表
	err = base_ar.SaveUser(tx, unitUserSaveDto.UserDto)
	if err != nil {
		return
	}
	// user_profile 表
	err = base_ar.SaveUserProfile(tx, unitUserSaveDto.UserProfileDto)
	if err != nil {
		return
	}
	// unit user表
	unitUserId, err = base_ar.UpsertUnitUser(tx, unitUserSaveDto.UnitUserDto.UnitUser, unitUserModel)
	if err != nil {
		return
	}

	// unit user_profile 表
	unitUserSaveDto.UnitUserProfileDto.Id = unitUserId
	err = base_ar.UpsertUnitUserProfile(tx, unitUserSaveDto.UnitUserProfileDto.UnitUserProfile, unitUserProfileModel)
	if err != nil {
		return
	}

	return
}
func (s *User) doEditUser(
	tx *gorm.DB,
	baseParamDto dto.BaseParamDto,
	unitUserSaveDto *user_dto.UnitUserSaveDto,
	userData user_dto.UserAllDataDto,
	unitUserData unit_dto.UnitUserAllDataDto,
	unitUserModel itf.UnitUserItf,
	unitUserProfileModel itf.UserProfileItf,
	unitDeptModel itf.DeptItf,
	unitUserDeptModel itf.UserDeptItf,
	unitRoleModel itf.RoleItf,
	unitUserRoleModel itf.UserRoleItf,
) (unitUserId string, err error) {
	// user表、user_profile表，
	// 主体用户由用户自己去负责更新

	// unit user表
	unitUserId, err = base_ar.UpsertUnitUser(tx, unitUserSaveDto.UnitUserDto.UnitUser, unitUserModel)
	if err != nil {
		return
	}
	// unit user_profile表
	err = base_ar.UpsertUnitUserProfile(tx, unitUserSaveDto.UnitUserProfileDto.UnitUserProfile, unitUserProfileModel)

	return
}

// 校验请求数据
func (s *User) checkRequestData(baseParamDto dto.BaseParamDto, unitUserSaveDto *user_dto.UnitUserSaveDto) (err error) {
	unitUserSaveDto.UnitUserProfileDto.CardNum = helper.DeleteSpace(unitUserSaveDto.UnitUserProfileDto.CardNum)
	unitUserSaveDto.UserProfileDto.CardNum = unitUserSaveDto.UnitUserProfileDto.CardNum

	unitUserSaveDto.UnitUserProfileDto.Remark = strings.TrimSpace(unitUserSaveDto.UnitUserProfileDto.Remark)
	unitUserSaveDto.UserProfileDto.Remark = unitUserSaveDto.UnitUserProfileDto.Remark

	unitUserSaveDto.UnitUserProfileDto.EmergencyName = strings.TrimSpace(unitUserSaveDto.UnitUserProfileDto.EmergencyName)
	unitUserSaveDto.UserProfileDto.EmergencyName = unitUserSaveDto.UnitUserProfileDto.EmergencyName

	unitUserSaveDto.UnitUserProfileDto.EmergencyTel = helper.DeleteSpace(unitUserSaveDto.UnitUserProfileDto.EmergencyTel)
	unitUserSaveDto.UserProfileDto.EmergencyTel = unitUserSaveDto.UnitUserProfileDto.EmergencyTel

	unitUserSaveDto.UnitUserDto.Phone = helper.DeleteSpace(unitUserSaveDto.UnitUserDto.Phone)
	unitUserSaveDto.UserDto.Phone = unitUserSaveDto.UnitUserDto.Phone

	unitUserSaveDto.UnitUserProfileDto.Company = helper.DeleteSpace(unitUserSaveDto.UnitUserProfileDto.Company)
	unitUserSaveDto.UserProfileDto.Company = unitUserSaveDto.UnitUserProfileDto.Company

	unitUserSaveDto.UnitUserProfileDto.Occupation = helper.DeleteSpace(unitUserSaveDto.UnitUserProfileDto.Occupation)
	unitUserSaveDto.UserProfileDto.Occupation = unitUserSaveDto.UnitUserProfileDto.Occupation

	unitUserSaveDto.UnitUserProfileDto.Address = strings.TrimSpace(unitUserSaveDto.UnitUserProfileDto.Address)
	unitUserSaveDto.UserProfileDto.Address = unitUserSaveDto.UnitUserProfileDto.Address

	unitUserSaveDto.UnitUserProfileDto.GraduatedFrom = helper.DeleteSpace(unitUserSaveDto.UnitUserProfileDto.GraduatedFrom)
	unitUserSaveDto.UserProfileDto.GraduatedFrom = unitUserSaveDto.UnitUserProfileDto.GraduatedFrom

	unitUserSaveDto.UnitUserProfileDto.Schooling = helper.DeleteSpace(unitUserSaveDto.UnitUserProfileDto.Schooling)
	unitUserSaveDto.UserProfileDto.Schooling = unitUserSaveDto.UnitUserProfileDto.Schooling

	unitUserSaveDto.UnitUserProfileDto.DegreeNumber = helper.DeleteSpace(unitUserSaveDto.UnitUserProfileDto.DegreeNumber)
	unitUserSaveDto.UserProfileDto.DegreeNumber = unitUserSaveDto.UnitUserProfileDto.DegreeNumber

	unitUserSaveDto.UnitUserProfileDto.Professional = helper.DeleteSpace(unitUserSaveDto.UnitUserProfileDto.Professional)
	unitUserSaveDto.UserProfileDto.Professional = unitUserSaveDto.UnitUserProfileDto.Professional

	unitUserSaveDto.UnitUserDto.UnitId = helper.DeleteSpace(unitUserSaveDto.UnitUserDto.UnitId)
	unitUserSaveDto.DeptId = helper.DeleteSpace(unitUserSaveDto.DeptId)
	// unitUserSaveDto.RoleId = unitUserSaveDto.RoleId

	if unitUserSaveDto.UnitUserDto.UnitId == "" {
		err = errors.New("请选择组织单位")
		return
	}
	if unitUserSaveDto.DeptId == "" {
		err = errors.New("请输入部门")
		return
	}
	if len(unitUserSaveDto.RoleId) == 0 {
		err = errors.New("请选择角色")
		return
	}

	if !helper.IsCellPhone(unitUserSaveDto.UnitUserDto.Phone) {
		err = errors.New("手机号码格式错误")
		return
	}

	unitUserSaveDto.UnitUserDto.Name = helper.DeleteSpace(unitUserSaveDto.UnitUserDto.Name)
	unitUserSaveDto.UserDto.Name = unitUserSaveDto.UnitUserDto.Name
	if unitUserSaveDto.UnitUserDto.Name == "" {
		err = errors.New("用户名称不能为空")
		return
	}

	unitUserSaveDto.UnitUserProfileDto.Email = helper.DeleteSpace(unitUserSaveDto.UnitUserProfileDto.Email)
	unitUserSaveDto.UserDto.Email = unitUserSaveDto.UnitUserProfileDto.Email
	if !helper.IsEmail(unitUserSaveDto.UserDto.Email) {
		err = errors.New("邮箱格式错误")
		return
	}

	return
}

// 校验用户与组织单位用户数据匹配
func (s *User) checkAndSetUserData(unitUserSaveDto *user_dto.UnitUserSaveDto, unitUserModel itf.UnitUserItf, unitUserProfileModel itf.UserProfileItf) (userData user_dto.UserAllDataDto, unitUserData unit_dto.UnitUserAllDataDto, err error) {
	isAddUnitUser := unitUserSaveDto.UnitUserDto.Id == ""

	if isAddUnitUser {
		if unitUserSaveDto.UserDto.Password == "" {
			err = errors.New("密码不能为空")
			return
		}
		if err = helper.CheckPasswordRule(unitUserSaveDto.UserDto.Password); err != nil {
			return
		}
	}
	tmpPwd, err1 := helper.GenerateCryptPassword(unitUserSaveDto.UserDto.Password)
	unitUserSaveDto.UserDto.Password = helper.Ternary(isAddUnitUser, tmpPwd, "")
	unitUserSaveDto.UserProfileDto.Deleted = 0
	unitUserSaveDto.UserProfileDto.Status = helper.Ternary(isAddUnitUser, models.USER_PROFILE_NORMAL, userData.Status)
	if err1 != nil {
		err = err1
		return
	}

	phone := unitUserSaveDto.UnitUserDto.Phone
	unitId := unitUserSaveDto.UnitUserDto.UnitId
	userId := &unitUserSaveDto.UnitUserDto.UserId
	if phone == "" {
		err = errors.New("checkAndSetUserData：手机号码不能为空")
		return
	}
	if unitId == "" {
		err = errors.New("checkAndSetUserData：组织单位id不能为空")
		return
	}
	// 1. 获取用户数据
	userAllInfo, err1 := base_ar.GetUserByIdOrPhone(unitUserSaveDto.UserDto.Id, phone)
	unitUserAllInfo, err2 := base_ar.GetUnitUserByIdOrPhone(unitUserSaveDto.UnitUserDto.Id, phone, unitId, unitUserModel, unitUserProfileModel)
	if err1 != nil || err2 != nil {
		err = helper.Ternary(err1 != nil, err1, err2)
		return
	}

	userAllInfoLen := len(userAllInfo)
	unitUserAllInfoLen := len(unitUserAllInfo)
	if userAllInfoLen == 1 && unitUserAllInfoLen == 1 {
		if userAllInfo[0].Id != unitUserAllInfo[0].UserId || (*userId != "" && *userId != userAllInfo[0].Id) {
			err = fmt.Errorf("用户（%s, %s）数据关联错误", phone, *userId)
			return
		}
		*userId = userAllInfo[0].Id
		userData = userAllInfo[0]
		unitUserData = unitUserAllInfo[0]
	} else {
		if userAllInfoLen > 1 || unitUserAllInfoLen > 1 {
			tmpTableName := helper.Ternary(userAllInfoLen > 1, (&models.User{}).TableName(), unitUserModel.TableName())
			err = fmt.Errorf("%s表存在多个相同手机号（%s），请核对数据", tmpTableName, phone)
			return
		}
		// else if userAllInfoLen != unitUserAllInfoLen {
		// 	err = fmt.Errorf("用户表数据条数异常，请核对数据")
		// 	return
		// }
	}

	// 2.页面数据转换, user表unit user表数据一起提交的，这里转换
	if isAddUnitUser {
		if unitUserAllInfoLen > 0 {
			err = errors.New("用户已存在，新增用户失败")
			return
		}
	} else {
		if unitUserAllInfoLen == 0 {
			err = errors.New("用户不存在，编辑用户失败")
			return
		}
	}
	// 已存在主体用户
	if userData.User.Id != "" && userData.Status == base_model.UNIT_USER_PROFILE_CANCLED {
		// 主体用户已注销-编辑同一个手机号
		if !isAddUnitUser {
			err = errors.New("主体用户已注销，请勿编辑")
			return
		}
		// 主体用户已注销-新增同一个手机号
		// pass
	} else if userData.User.Id != "" && userData.Status != base_model.UNIT_USER_PROFILE_CANCLED {
		// 主体用户未注销-新增同一个手机号
		if isAddUnitUser {
			err = fmt.Errorf("用户（%s）已存在，请勿重复新增", unitUserSaveDto.UnitUserDto.Phone)
			return
		}
		// 主体用户未注销-编辑同一个手机号
		// pass
	}

	// 3.处理主体用户和组织单位用户id 关系
	if isAddUnitUser {
		uuid := ""
		uuid, err = helper.GetUuid()
		unitUserSaveDto.UnitUserDto.UserId = *userId
		unitUserSaveDto.UnitUserDto.Id = uuid
		unitUserSaveDto.UnitUserDto.UnitUser.Id = uuid
		unitUserSaveDto.UnitUserProfileDto.Id = uuid
		unitUserSaveDto.UnitUserProfileDto.UnitUserProfile.Id = uuid
		if err != nil {
			return
		}

		if userData.User.Id == "" || userData.Status == base_model.UNIT_USER_PROFILE_CANCLED {
			err = s.setNewUserId(unitUserSaveDto)
			if err != nil {
				return
			}
		}
	}

	unitUserSaveDto.UnitUserDeptDto, err = s.generateUserDeptData(unitUserSaveDto.UnitUserDto.Id, unitUserSaveDto.DeptId)
	if err != nil {
		return
	}
	unitUserSaveDto.UnitUserRoleDto, err = s.generateUserRoleData(unitUserSaveDto.UnitUserDto.Id, unitUserSaveDto.RoleId)
	if err != nil {
		return
	}
	return
}

// 校验组织结构
func (s *User) CheckOrgStructure(unitUserId string, unitId string, deptId string, roleId []string, unitUserModel itf.UnitUserItf, unitDeptModel itf.DeptItf, unitRoleModel itf.RoleItf) (err error) {
	result := struct {
		ExistUnit int `gorm:"column:exist_unit"`
		ExistDept int `gorm:"column:exist_dept"`
		ExistRole int `gorm:"column:exist_role"`
	}{}

	err = base_ar.CheckOrgStructure(&result, unitUserId, unitId, deptId, roleId, unitUserModel, unitDeptModel, unitRoleModel)
	if err != nil {
		return
	}
	if result.ExistUnit == 0 {
		err = fmt.Errorf("组织单位（%s）不存在", unitId)
		return
	}
	if result.ExistDept == 0 {
		err = fmt.Errorf("部门（%s）不存在", deptId)
		return
	}
	if result.ExistRole == 0 {
		err = fmt.Errorf("角色（%s）不存在", roleId)
		return
	}
	return nil
}

// 设置主体用户id
func (s *User) setNewUserId(unitUserSaveDto *user_dto.UnitUserSaveDto) (err error) {
	newUserId, err := helper.GetUuid()
	if err != nil {
		return
	}
	unitUserSaveDto.UserDto.Id = newUserId
	unitUserSaveDto.UserDto.User.Id = newUserId
	unitUserSaveDto.UserProfileDto.Id = newUserId
	unitUserSaveDto.UserProfileDto.UserProfile.Id = newUserId
	unitUserSaveDto.UnitUserDto.UserId = newUserId
	unitUserSaveDto.UnitUserDto.UnitUser.UserId = newUserId
	return
}
func (s *User) generateUserDeptData(unitUserId, deptId string) (data user_dto.UnitUserDeptDto, err error) {
	uuid := ""
	uuid, err = helper.GetUuid()
	if err != nil {
		return
	}
	data = user_dto.UnitUserDeptDto{
		UnitUserDept: base_model.UnitUserDept{
			Id:      uuid,
			UserId:  unitUserId,
			DeptId:  deptId,
			Deleted: 0,
		},
	}
	return
}
func (s *User) generateUserRoleData(unitUserId string, roleIds []string) (data []user_dto.UnitUserRoleDto, err error) {

	data = make([]user_dto.UnitUserRoleDto, 0)
	for _, v := range roleIds {
		uuid := ""
		uuid, err = helper.GetUuid()
		if err != nil {
			return
		}
		item := user_dto.UnitUserRoleDto{
			UnitUserRole: base_model.UnitUserRole{
				Id:      uuid,
				UserId:  unitUserId,
				RoleId:  v,
				Deleted: 0,
			},
		}
		data = append(data, item)
	}
	return
}

// 删除组织单位用户
func (s *User) DelUnitUser(baseParamDto dto.BaseParamDto, ids []string) error {

	switch baseParamDto.ModuleName {
	case "admin_plat":
		return base_ar.DelUnitUser(ids, &models.PlatUser{}, &models.PlatUserProfile{})
	case "admin_mchnt":
		return base_ar.DelUnitUser(ids, &models.MchntUser{}, &models.MchntUserProfile{})
	default:
		return errors.New("DelUnitUser:模块名称错误")
	}
}

func (s *User) DelUnitUserForAdminPlat(baseParamDto dto.BaseParamDto, ids []string) error {
	if true {
		return errors.New("平台没有操作权限")
	}
	return nil
}
