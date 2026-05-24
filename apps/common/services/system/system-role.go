package system

import (
	"WenBeego/apps/common/dto"
	"WenBeego/apps/common/dto/page_dto"
	"WenBeego/apps/common/dto/role_dto"
	"WenBeego/apps/common/global"
	"WenBeego/apps/common/helper"
	"WenBeego/apps/common/middleware/business_store"
	"WenBeego/apps/common/models"
	"WenBeego/apps/common/models/base_model"
	"WenBeego/apps/common/models/itf"
	"WenBeego/apps/common/models_ar"
	"WenBeego/apps/common/models_ar/base_ar"
	"errors"

	"gorm.io/gorm"
)

type Role struct {
	PlatMenuViewAr models_ar.PlatMenuViewAr
}

// 获取角色列表
func (s *Role) GetUnitRoleList(RoleDto page_dto.SystemRoleListReqDto) (resultDto dto.RespDataListDto, err error) {
	data := make([]base_model.UnitRole, 0)
	var count int64 = 0

	if res, err1 := helper.CheckUserHasUnit(RoleDto.ModuleName, RoleDto.UserId, RoleDto.SelectUnitIds); !res {
		err = helper.Ternary(err1 != nil, err1, errors.New("GetUnitRoleList：用户没有组织单位权限"))
		return
	}

	switch RoleDto.ModuleName {
	case "admin_plat":
		data, count, err = base_ar.GetUnitRoleList(RoleDto, &models.Plat{}, &models.PlatRole{}, &models.PlatRoleClassify{})
	case "admin_mchnt":
		data, count, err = base_ar.GetUnitRoleList(RoleDto, &models.Mchnt{}, &models.MchntRole{}, &models.MchntRoleClassify{})
	default:
		err = errors.New("GetUnitList:模块名称错误")
	}
	if err != nil {
		return
	}
	resultDto = dto.RespDataListDto{Total: count, List: data, PageSize: RoleDto.PageSize, CurrentPage: RoleDto.CurrentPage}
	return
}
func (s *Role) GetUnitRoleListForAdminPlat(RoleDto page_dto.SystemRoleListReqDto) (resultDto dto.RespDataListDto, err error) {
	data := make([]base_model.UnitRole, 0)
	var count int64 = 0

	// if res, err1 := helper.CheckUserHasUnit(RoleDto.ModuleName, RoleDto.UserId, RoleDto.SelectUnitIds); !res {
	// 	err = helper.Ternary(err1 != nil, err1, errors.New("GetUnitRoleList：用户没有组织单位权限"))
	// 	return
	// }

	data, count, err = base_ar.GetUnitRoleList(RoleDto, &models.Mchnt{}, &models.MchntRole{}, &models.MchntRoleClassify{})
	if err != nil {
		return
	}
	resultDto = dto.RespDataListDto{Total: count, List: data, PageSize: RoleDto.PageSize, CurrentPage: RoleDto.CurrentPage}
	return
}

// 保存角色列表
func (s *Role) SaveUnitRole(baseParamDto dto.BaseParamDto, roleDto role_dto.UnitRoleDto) (id string, err error) {
	roleDto.RoleName = helper.DeleteSpace(roleDto.RoleName)
	roleDto.RoleClassifyName = helper.DeleteSpace(roleDto.RoleClassifyName)
	if roleDto.RoleName == "" {
		err = errors.New("请输入角色名称")
		return
	}
	if roleDto.RoleClassifyName == "" {
		err = errors.New("请输入角色分类名称")
		return
	}

	if res, err1 := helper.CheckUserHasUnit(baseParamDto.ModuleName, baseParamDto.UserId, []string{roleDto.UnitId}); !res {
		err = helper.Ternary(err1 != nil, err1, errors.New("SaveUnitRole：用户没有组织单位权限"))
		return
	}

	var classifyData base_model.UnitRoleClassify
	if roleDto.Id == "" {
		var err1 error
		roleDto.Id, err1 = helper.GetUuid()
		roleDto.CreatedBy = baseParamDto.UnitUserId

		uuid := ""
		uuid, err = helper.GetUuid()
		if err != nil || err1 != nil {
			err = helper.Ternary(err != nil, err, err1)
			return
		}
		classifyData = base_model.UnitRoleClassify{
			Id:      uuid,
			RoleId:  roleDto.Id,
			UnitId:  roleDto.UnitId,
			Name:    roleDto.RoleClassifyName,
			Deleted: 0,
		}
	} else {
		switch baseParamDto.ModuleName {
		case "admin_plat":
			classifyData, err = base_ar.GetRoleClassifyByRoleId(roleDto.Id, &models.PlatRoleClassify{})
		case "admin_mchnt":
			classifyData, err = base_ar.GetRoleClassifyByRoleId(roleDto.Id, &models.MchntRoleClassify{})
		default:
			err = errors.New("模块名称错误")
		}

		roleDto.UpdatedBy = baseParamDto.UnitUserId
		classifyData.Name = roleDto.RoleClassifyName
		classifyData.UnitId = roleDto.UnitId
	}
	if err != nil {
		return
	}

	err = global.WriteDb.Transaction(func(tx *gorm.DB) (err error) {
		switch baseParamDto.ModuleName {
		case "admin_plat":
			id, err = base_ar.SaveUnitRole(tx, roleDto, &models.PlatRole{})
			base_ar.InsertUserRoleClassifies[*models.PlatRoleClassify](tx, classifyData)
		case "admin_mchnt":
			id, err = base_ar.SaveUnitRole(tx, roleDto, &models.MchntRole{})
			base_ar.InsertUserRoleClassifies[*models.MchntRoleClassify](tx, classifyData)
		default:
			err = errors.New("模块名称错误")
		}
		return
	})

	// 清空用户权限认证缓存
	business_store.ClearAumid()

	return
}
func (s *Role) SaveUnitRoleForAdminPlat(baseParamDto dto.BaseParamDto, roleDto role_dto.UnitRoleDto) (id string, err error) {
	if true {
		err = errors.New("平台没有操作权限")
		return
	}
	return
}

// 删除角色列表
func (s *Role) DelUnitRole(baseParamDto dto.BaseParamDto, roleDto role_dto.UnitRoleDto) (err error) {
	timestamp := helper.GetTimestamp()
	updateData := base_model.UnitRole{}
	updateData.Id = roleDto.Id
	updateData.UpdatedAt = timestamp
	updateData.Deleted = 1

	switch baseParamDto.ModuleName {
	case "admin_plat":
		err = base_ar.DelUnitRole(updateData, &models.PlatRole{})
	case "admin_mchnt":
		err = base_ar.DelUnitRole(updateData, &models.MchntRole{})
	default:
		err = errors.New("模块名称错误")
	}
	return
}
func (s *Role) DelUnitRoleForAdminPlat(baseParamDto dto.BaseParamDto, roleDto role_dto.UnitRoleDto) (err error) {
	if true {
		err = errors.New("平台没有操作权限")
		return
	}
	return
}

func (s *Role) GetRoleMenu(baseParamDto dto.BaseParamDto, selectUnitIds []string) (data interface{}, err error) {
	if res, err1 := helper.CheckUserHasUnit(baseParamDto.ModuleName, baseParamDto.UserId, selectUnitIds); !res {
		err = helper.Ternary(err1 != nil, err1, errors.New("GetRoleMenu：用户没有组织单位权限"))
		return
	}

	var dataList interface{}
	switch baseParamDto.ModuleName {
	case "admin_plat":
		// dataList, err = base_ar.GetRoleMenu(selectUnitIds, &models.PlatMenu{}, &models.PlatMenuMap{})
		dataList, err = s.PlatMenuViewAr.GetRoleMenu(selectUnitIds, models.PlatMenuView{}, models.PlatMenuMap{})
	case "admin_mchnt":
		dataList, err = base_ar.GetRoleMenu(selectUnitIds, &models.MchntMenu{}, &models.MchntMenuMap{})
	default:
		err = errors.New("GetRoleMenu：模块名称错误")
	}
	if err != nil {
		return
	}
	data = struct {
		List interface{} `json:"list"`
	}{List: dataList}

	return
}
func (s *Role) GetRoleMenuForAdminPlat(baseParamDto dto.BaseParamDto, selectUnitIds []string) (data interface{}, err error) {
	// if res, err1 := helper.CheckUserHasUnit(baseParamDto.ModuleName, baseParamDto.UserId, selectUnitIds); !res {
	// 	err = helper.Ternary(err1 != nil, err1, errors.New("GetRoleMenu：用户没有组织单位权限"))
	// 	return
	// }

	var dataList interface{}
	dataList, err = base_ar.GetRoleMenu(selectUnitIds, &models.MchntMenu{}, &models.MchntMenuMap{})
	if err != nil {
		return
	}
	data = struct {
		List interface{} `json:"list"`
	}{List: dataList}

	return
}
func (s *Role) GetRoleMenuIds(baseParamDto dto.BaseParamDto, roleId string) (data interface{}, err error) {

	dataList := make([]base_model.UnitRoleMenu, 0)

	switch baseParamDto.ModuleName {
	case "admin_plat":
		// dataList, err = base_ar.GetRoleMenuIds(roleId, &models.PlatMenu{}, &models.PlatMenuMap{}, &models.PlatRoleMenu{})
		dataList, err = s.PlatMenuViewAr.GetRoleMenuIds(baseParamDto, roleId, models.PlatMenuView{}, models.PlatMenuMapView{}, models.PlatRoleMenu{})
	case "admin_mchnt":
		dataList, err = base_ar.GetRoleMenuIds(roleId, &models.MchntMenu{}, &models.MchntMenuMap{}, &models.MchntRoleMenu{})
	default:
		err = errors.New("GetRoleMenuIds：模块名称错误")
	}
	if err != nil {
		return
	}
	tmpData := make([]string, 0)
	for _, v := range dataList {
		tmpData = append(tmpData, v.MenuId)
	}
	data = struct {
		List []string `json:"list"`
	}{List: tmpData}

	return
}
func (s *Role) GetRoleMenuIdsForAdminPlat(baseParamDto dto.BaseParamDto, roleId string) (data interface{}, err error) {

	dataList := make([]base_model.UnitRoleMenu, 0)
	dataList, err = base_ar.GetRoleMenuIds(roleId, &models.MchntMenu{}, &models.MchntMenuMap{}, &models.MchntRoleMenu{})
	if err != nil {
		return
	}
	tmpData := make([]string, 0)
	for _, v := range dataList {
		tmpData = append(tmpData, v.MenuId)
	}
	data = struct {
		List []string `json:"list"`
	}{List: tmpData}

	return
}

// 保存角色所选菜单
func (s *Role) RoleMenuSave(baseParamDto dto.BaseParamDto, roleMenuSaveDto role_dto.RoleMenuSaveDto) (err error) {
	roleId := roleMenuSaveDto.RoleId
	menuIds := roleMenuSaveDto.MenuIds

	if roleId == "" {
		return errors.New("RoleMenuSave：角色ID不能为空")
	}

	err = global.WriteDb.Transaction(func(tx *gorm.DB) (err error) {
		switch baseParamDto.ModuleName {
		case "admin_plat":
			roleClassify, err1 := base_ar.GetRoleClassifyByRoleId(roleId, &models.PlatRoleClassify{})
			if err1 == nil && roleClassify.Id != "" && roleClassify.Name == "admin" {
				return errors.New("RoleMenuSave：系统管理员角色不能修改菜单")
			}
			err = s.doRoleMenuSave(tx, roleId, menuIds, &models.PlatRoleMenu{})
		case "admin_mchnt":
			roleClassify, err1 := base_ar.GetRoleClassifyByRoleId(roleId, &models.MchntRoleClassify{})
			if err1 == nil && roleClassify.Id != "" && roleClassify.Name == "admin" {
				return errors.New("RoleMenuSave：系统管理员的角色不能修改菜单")
			}
			err = s.doRoleMenuSave(tx, roleId, menuIds, &models.MchntRoleMenu{})
		default:
			err = errors.New("RoleMenuSave：模块名称错误")
		}
		return
	})
	// 清空用户权限认证缓存
	business_store.ClearAumid()
	return
}
func (s *Role) RoleMenuSaveForAdminPlat(baseParamDto dto.BaseParamDto, roleMenuSaveDto role_dto.RoleMenuSaveDto) (err error) {
	if true {
		err = errors.New("平台没有操作权限")
		return
	}
	return
}

// 保存角色所选菜单
func (s *Role) doRoleMenuSave(tx *gorm.DB, roleId string, menuIds []string, unitRoleMenuModel itf.RoleMenuItf) (err error) {
	err = base_ar.DeleteRoleMenuByRole(tx, roleId, unitRoleMenuModel)
	if err != nil {
		return
	}
	if len(menuIds) > 0 {
		err = base_ar.RoleMenuSave(tx, roleId, menuIds, unitRoleMenuModel)
	}
	return err
}
