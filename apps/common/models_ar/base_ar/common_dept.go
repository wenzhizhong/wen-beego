package base_ar

import (
	"WenBeego/apps/common/dto_vo/dept_dto"
	"WenBeego/apps/common/dto_vo/page_dto"
	"WenBeego/apps/common/global"
	"WenBeego/apps/common/helper"
	"WenBeego/apps/common/models"
	"WenBeego/apps/common/models/base_model"
	"WenBeego/apps/common/models/itf"
	"errors"
)

// 获取组织架构部门列表
func GetUnitDeptList[UnitModel itf.UnitItf, UnitDeptModel itf.DeptItf, UnitUserModel itf.UnitUserItf, UnitUserProfileModel itf.UserProfileItf](deptDto page_dto.SystemDeptListReqDto, unitModel UnitModel, unitDeptModel UnitDeptModel, unitUserModel UnitUserModel, unitUserProfileModel UnitUserProfileModel) (dataList []base_model.UnitDept, count int64, err error) {
	tableUnitName := unitModel.TableName()
	tableDeptName := unitDeptModel.TableName()
	tableUnitUserName := unitUserModel.TableName()
	tableUserProfileName := unitUserProfileModel.TableName()

	subQuery := global.GetReadDb().
		Model(unitUserModel).
		Select(tableUnitUserName+".id AS principal_id, "+tableUnitUserName+".name AS principal, "+tableUnitUserName+".phone, "+tableUserProfileName+".email").
		Joins("INNER JOIN "+tableUserProfileName+" ON "+tableUserProfileName+".id = "+tableUnitUserName+".id").
		Where(tableUnitUserName+".deleted = 0").
		Where(tableUserProfileName+".deleted = 0").
		Where(tableUserProfileName+".status = ?", base_model.UNIT_USER_PROFILE_NORMAL)
	if len(deptDto.SelectUnitIds) > 0 {
		subQuery = subQuery.Where(tableUnitUserName+".unit_id in ?", deptDto.SelectUnitIds)
	}

	query := global.GetReadDb().
		Model(unitDeptModel).
		Joins("INNER JOIN " + tableUnitName + " ON " + tableUnitName + ".id = " + tableDeptName + ".unit_id").
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
	err = query.Select(tableDeptName+".*,"+tableUnitName+".name as unit_name, t.principal, t.phone, t.email").
		Joins("LEFT JOIN (?) AS t ON t.principal_id = "+tableDeptName+".principal_id", subQuery).
		Order(tableDeptName + ".sort").
		Find(&dataList).Error
	return
}

// 获取组织架构树
func GetUnitDeptTree[UnitDeptModel itf.DeptItf](selectUnitIds []string, unitDeptModel UnitDeptModel) (data []base_model.UnitDept, err error) {
	data = make([]base_model.UnitDept, 0)
	if len(selectUnitIds) == 0 {
		return data, nil
	}
	tableDeptName := unitDeptModel.TableName()
	err = global.GetReadDb().
		Model(unitDeptModel).
		Select(tableDeptName+".id,"+tableDeptName+".pid,"+tableDeptName+".name").
		Where("unit_id in (?)", selectUnitIds).
		Where("deleted = 0").
		Order(tableDeptName + ".sort").
		Find(&data).Error
	return
}

// 获取可用的部门负责人
func GetUnitDeptPrincipal[UnitUserModel itf.UnitUserItf, UnitUserProfileModel itf.UserProfileItf](deptPrincipalDto page_dto.SystemDeptPrincipalReqDto, unitUserModel UnitUserModel, unitUserProfileModel UnitUserProfileModel) (data interface{}, count int64, err error) {
	type dataStruct struct {
		Id   string `json:"id"`
		Name string `json:"name"`
	}

	data = make([]dataStruct, 0)
	if len(deptPrincipalDto.SelectUnitIds) == 0 {
		return data, 0, nil
	}
	tableUserName := (&models.User{}).TableName()
	tableUserProfileName := (&models.UserProfile{}).TableName()
	tableUnitUserName := unitUserModel.TableName()
	tableUnitUserProfileName := unitUserProfileModel.TableName()

	tableStruct := struct {
		User            string
		UserProfile     string
		UnitUser        string
		UnitUserProfile string
	}{
		User:            tableUserName,
		UserProfile:     tableUserProfileName,
		UnitUser:        tableUnitUserName,
		UnitUserProfile: tableUnitUserProfileName,
	}
	selectTpl := `"{{.UnitUser}}".id, 
		CASE WHEN {{.UnitUser}}.name!='' THEN {{.UnitUser}}.name ELSE "{{.User}}".name END AS name
	`
	joinTpl := `
		INNER JOIN "{{.User}}" ON "{{.User}}".id = "{{.UnitUser}}".user_id
		INNER JOIN {{.UserProfile}} ON {{.UserProfile}}.id = "{{.UnitUser}}".user_id
		INNER JOIN {{.UnitUserProfile}} ON {{.UnitUserProfile}}.id = "{{.UnitUser}}".id
	`
	selectStr, err := helper.ParseStringTpl(selectTpl, tableStruct)
	joinStr, err2 := helper.ParseStringTpl(joinTpl, tableStruct)

	if err != nil {
		return data, 0, err
	}
	if err2 != nil {
		return data, 0, err2
	}

	query := global.GetReadDb().
		Model(unitUserModel).
		Select(selectStr).
		Joins(joinStr).
		Where(tableUnitUserName+".unit_id IN (?)", deptPrincipalDto.SelectUnitIds).
		Where(tableUnitUserName+".deleted = 0").
		Where("\""+tableUserProfileName+"\".status = ?", base_model.UNIT_USER_PROFILE_NORMAL).
		Where("\""+tableUserProfileName+"\".deleted = ?", 0).
		Where("\""+tableUnitUserProfileName+"\".status = ?", base_model.UNIT_USER_PROFILE_NORMAL).
		Where("\""+tableUnitUserProfileName+"\".deleted = ?", 0)

	if deptPrincipalDto.Keyword != "" {
		query = query.Where("\""+tableUserName+"\".name like ?", "%"+deptPrincipalDto.Keyword+"%")
	}

	err = query.Count(&count).Error
	if err != nil || count == 0 {
		return data, 0, err
	}

	err = query.Find(&data).Error
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
			Select("*").
			Where("id = ?", unitDeptDto.Id).
			Updates(unitDeptDto).Error
	}
	return unitDeptDto.Id, err
}

// 删除组织架构
func DelUnitDept[UnitDeptModel itf.DeptItf](unitDeptData base_model.UnitDept, unitDeptModel UnitDeptModel) error {
	if unitDeptData.Id == "" {
		return errors.New("DelUnitDept: 参数id不能为空")
	}
	return global.GetWriteDb().
		Model(unitDeptModel).
		Where("id = ?", unitDeptData.Id).
		Updates(unitDeptData).Error
}
