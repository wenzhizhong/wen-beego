package system

import (
	"WenBeego/apps/common/dto"
	"WenBeego/apps/common/dto/page_dto"
	"WenBeego/apps/common/dto/role_dto"
	"WenBeego/apps/common/global"
	"WenBeego/apps/common/helper"
	"WenBeego/apps/common/models"
	"WenBeego/apps/common/models/base_model"
	"WenBeego/apps/common/models/itf"
	"WenBeego/apps/common/models_ar/base_ar"
	"errors"
	"regexp"

	"gorm.io/gorm"
)

type Role struct {
}

// 获取角色列表
func (s *Role) GetUnitRoleList(RoleDto page_dto.SystemRoleListReqDto) (resultDto dto.RespDataListDto, err error) {
	data := make([]base_model.UnitRole, 0)
	var count int64 = 0

	switch RoleDto.ModuleName {
	case "admin_plat":
		data, count, err = base_ar.GetUnitRoleList(RoleDto, &models.Plat{}, &models.PlatRole{}, &models.PlatRoleClassify{})
	case "mchnt_plat":
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

// 保存角色列表
func (s *Role) SaveUnitRole(baseParamDto dto.BaseParamDto, RoleDto role_dto.UnitRoleDto) (id string, err error) {
	RoleDto.UnitId = baseParamDto.UnitId
	re := regexp.MustCompile(`\s+`)
	RoleDto.RoleClassifyName = re.ReplaceAllString(RoleDto.RoleClassifyName, "")

	var classifyData base_model.UnitRoleClassify
	if RoleDto.Id == "" {
		RoleDto.Id, _ = helper.GetUuid()
		RoleDto.CreatedBy = baseParamDto.UnitUserId

		uuid, _ := helper.GetUuid()
		classifyData = base_model.UnitRoleClassify{
			Id:      uuid,
			RoleId:  RoleDto.Id,
			UnitId:  RoleDto.UnitId,
			Name:    RoleDto.RoleClassifyName,
			Deleted: 0,
		}
	} else {
		switch baseParamDto.ModuleName {
		case "admin_plat":
			classifyData, err = base_ar.GetRoleClassifyByRoleId(RoleDto.Id, &models.PlatRoleClassify{})
		case "mchnt_plat":
			classifyData, err = base_ar.GetRoleClassifyByRoleId(RoleDto.Id, &models.MchntRoleClassify{})
		default:
			err = errors.New("模块名称错误")
		}

		RoleDto.UpdatedBy = baseParamDto.UnitUserId
		classifyData.Name = RoleDto.RoleClassifyName
		classifyData.UnitId = RoleDto.UnitId
	}
	if err != nil {
		return
	}

	err = global.WriteDb.Transaction(func(tx *gorm.DB) (err error) {
		switch baseParamDto.ModuleName {
		case "admin_plat":
			id, err = base_ar.SaveUnitRole(tx, RoleDto, &models.PlatRole{})
			base_ar.InsertUserRoleClassifies[*models.PlatRoleClassify](tx, classifyData)
		case "mchnt_plat":
			id, err = base_ar.SaveUnitRole(tx, RoleDto, &models.MchntRole{})
			base_ar.InsertUserRoleClassifies[*models.MchntRoleClassify](tx, classifyData)
		default:
			err = errors.New("模块名称错误")
		}
		return
	})

	return
}

// 删除角色列表
func (s *Role) DelUnitRole(baseParamDto dto.BaseParamDto, RoleDto role_dto.UnitRoleDto) (err error) {
	timestamp := helper.GetTimestamp()
	updateData := base_model.UnitRole{}
	updateData.Id = RoleDto.Id
	updateData.UpdatedAt = timestamp
	updateData.Deleted = 1

	switch baseParamDto.ModuleName {
	case "admin_plat":
		err = base_ar.DelUnitRole(updateData, &models.PlatRole{})
	case "mchnt_plat":
		err = base_ar.DelUnitRole(updateData, &models.MchntRole{})
	default:
		err = errors.New("模块名称错误")
	}
	return
}

func (s *Role) GetRoleMenu(baseParamDto dto.BaseParamDto, selectUnitIds []string) (data interface{}, err error) {
	var dataList interface{}
	switch baseParamDto.ModuleName {
	case "admin_plat":
		dataList, err = base_ar.GetRoleMenu(selectUnitIds, &models.PlatMenu{})
	case "mchnt_plat":
		dataList, err = base_ar.GetRoleMenu(selectUnitIds, &models.MchntMenu{})
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
func (s *Role) GetRoleMenuIds(baseParamDto dto.BaseParamDto, roleId string) (data interface{}, err error) {

	dataList := make([]base_model.UnitRoleMenu, 0)

	switch baseParamDto.ModuleName {
	case "admin_plat":
		dataList, err = base_ar.GetRoleMenuIds(roleId, &models.PlatMenu{}, &models.PlatRoleMenu{})
	case "mchnt_plat":
		dataList, err = base_ar.GetRoleMenuIds(roleId, &models.MchntMenu{}, &models.MchntRoleMenu{})
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
		case "mchnt_plat":
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
