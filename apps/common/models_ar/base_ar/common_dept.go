package base_ar

import (
	"WenBeego/apps/common/dto/dept_dto"
	"WenBeego/apps/common/dto/page_dto"
	"WenBeego/apps/common/global"
	"WenBeego/apps/common/helper"
	"WenBeego/apps/common/models/base_model"
	"WenBeego/apps/common/models/itf"
)

// 获取组织架构部门列表
func GetUnitDeptList[UnitModel itf.UnitItf, UnitDeptModel itf.DeptItf](deptDto page_dto.SystemDeptListReqDto, unitModel UnitModel, unitDeptModel UnitDeptModel) (dataList []base_model.UnitDept, count int64, err error) {
	tableUnitName := unitModel.TableName()
	tableDeptName := unitDeptModel.TableName()

	query := global.GetReadDb().
		Model(unitDeptModel).
		Joins("inner join " + tableUnitName + " on " + tableUnitName + ".id = " + tableDeptName + ".unit_id").
		Where(tableDeptName + ".deleted=0").
		Where(tableUnitName + ".deleted=0")
	if deptDto.Name != "" {
		query = query.Where(tableDeptName+".name like ?", "%"+deptDto.Name+"%")
	}
	if len(deptDto.SelectUnitIds) > 0 {
		query = query.Where(tableDeptName+".unit_id in ?", deptDto.SelectUnitIds)
	}

	err = query.Select(tableDeptName + ".id").Count(&count).Error
	if err != nil {
		return make([]base_model.UnitDept, 0), 0, nil
	}
	if count == 0 {
		return make([]base_model.UnitDept, 0), 0, nil
	}
	err = query.Select(tableDeptName + ".*").
		Order(tableDeptName + ".sort").
		Find(&dataList).Error
	return
}

// 保存部门
func SaveUnitDept[UnitDeptModel itf.DeptItf](unitDeptDto dept_dto.UnitDeptDto, unitDeptModel UnitDeptModel) (id string, err error) {
	if unitDeptDto.Id == "" {
		unitDeptDto.Id, err = helper.GetUuid()
		if err != nil {
			return
		}

		err = global.GetWriteDb().
			Model(unitDeptModel).
			Create(&unitDeptDto).Error
	} else {
		err = global.GetWriteDb().
			Model(unitDeptModel).
			Where("id = ?", unitDeptDto.Id).
			Updates(unitDeptDto).Error
	}
	return unitDeptDto.Id, err
}
