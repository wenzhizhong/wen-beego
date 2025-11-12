package system

import (
	"WenBeego/apps/common/dto"
	"WenBeego/apps/common/dto/dept_dto"
	"WenBeego/apps/common/dto/page_dto"
	"WenBeego/apps/common/helper"
	"WenBeego/apps/common/models"
	"WenBeego/apps/common/models/base_model"
	"WenBeego/apps/common/models_ar/base_ar"
	"errors"
)

type Dept struct {
}

// 获取组织架构列表
func (s *Dept) GetUnitDeptList(deptDto page_dto.SystemDeptListReqDto) (resultDto dto.RespDataListDto, err error) {
	data := make([]base_model.UnitDept, 0)
	var count int64 = 0

	switch deptDto.ModuleName {
	case "admin_plat":
		data, count, err = base_ar.GetUnitDeptList(deptDto, &models.Plat{}, &models.PlatDept{}, &models.PlatUser{}, &models.PlatUserProfile{})
	case "mchnt_plat":
		data, count, err = base_ar.GetUnitDeptList(deptDto, &models.Mchnt{}, &models.MchntDept{}, &models.MchntUser{}, &models.MchntUserProfile{})
	default:
		err = errors.New("GetUnitList:模块名称错误")
	}
	if err != nil {
		return
	}
	resultDto = dto.RespDataListDto{Total: count, List: data}
	return
}
func (s *Dept) GetUnitDeptTree(baseParamDto dto.BaseParamDto, selectUnitIds []string) (data interface{}, err error) {
	dataList := make([]base_model.UnitDept, 0)

	switch baseParamDto.ModuleName {
	case "admin_plat":
		dataList, err = base_ar.GetUnitDeptTree(selectUnitIds, &models.PlatDept{})
	case "mchnt_plat":
		dataList, err = base_ar.GetUnitDeptTree(selectUnitIds, &models.MchntDept{})
	default:
		err = errors.New("GetUnitDeptTree:模块名称错误")
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

// 获取可用组织架构负责人
func (s *Dept) GetUnitDeptPrincipal(baseParamDto dto.BaseParamDto, deptPrincipalDto page_dto.SystemDeptPrincipalReqDto) (resultDto interface{}, err error) {
	var count int64
	var data interface{}

	switch baseParamDto.ModuleName {
	case "admin_plat":
		data, count, err = base_ar.GetUnitDeptPrincipal(deptPrincipalDto, &models.PlatUser{}, &models.PlatUserProfile{})
	case "mchnt_plat":
		data, count, err = base_ar.GetUnitDeptPrincipal(deptPrincipalDto, &models.MchntUser{}, &models.MchntUserProfile{})
	default:
		err = errors.New("GetUnitDeptPrincipal:模块名称错误")
	}
	if err != nil {
		return
	}
	resultDto = dto.RespDataListDto{Total: count, List: data}
	return
}

// 保存组织架构列表
func (s *Dept) SaveUnitDept(baseParamDto dto.BaseParamDto, deptDto dept_dto.UnitDeptDto) (id string, err error) {
	deptDto.UnitId = baseParamDto.UnitId
	switch baseParamDto.ModuleName {
	case "admin_plat":
		id, err = base_ar.SaveUnitDept(deptDto, &models.PlatDept{})
	case "mchnt_plat":
		id, err = base_ar.SaveUnitDept(deptDto, &models.MchntDept{})
	default:
		err = errors.New("模块名称错误")
	}
	return
}

// 删除组织架构列表
func (s *Dept) DelUnitDept(baseParamDto dto.BaseParamDto, deptDto dept_dto.UnitDeptDto) (err error) {
	timestamp := helper.GetTimestamp()
	updateData := base_model.UnitDept{}
	updateData.Id = deptDto.Id
	updateData.UpdatedAt = timestamp
	updateData.DeletedAt = &timestamp
	updateData.Deleted = 1

	switch baseParamDto.ModuleName {
	case "admin_plat":
		err = base_ar.DelUnitDept(updateData, &models.PlatDept{})
	case "mchnt_plat":
		err = base_ar.DelUnitDept(updateData, &models.MchntDept{})
	default:
		err = errors.New("模块名称错误")
	}
	return
}
