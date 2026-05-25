package base_ar

import (
	"WenBeego/apps/common/dto_vo/page_dto"
	"WenBeego/apps/common/dto_vo/unit_dto"
	"WenBeego/apps/common/dto_vo/unit_vo"
	"WenBeego/apps/common/dto_vo/user_vo"
	"WenBeego/apps/common/global"
	"WenBeego/apps/common/helper"
	"WenBeego/apps/common/models"
	"WenBeego/apps/common/models/base_model"
	"WenBeego/apps/common/models/itf"
	"errors"
	"fmt"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// 通过id或者手机号，获取组织单位用户
func GetUnitUserByIdOrPhone[UnitUserModel itf.UnitUserItf, UnitUserProfileModel itf.UserProfileItf](id, phone, unitId string, unitUserModel UnitUserModel, unitUserProfileModel UnitUserProfileModel) (data []unit_dto.UnitUserAllDataDto, err error) {
	if id == "" && phone == "" {
		return nil, errors.New("GetUnitUserByIdOrPhone(): Id 和 Phone 不能同时为空")
	}
	if unitId == "" {
		return nil, errors.New("GetUnitUserByIdOrPhone(): UnitId 不能为空")
	}
	tableUnitUserName := unitUserModel.TableName()
	tableUnitUserProfileName := unitUserProfileModel.TableName()

	query := global.GetReadDb().
		Model(unitUserModel).
		Select("*").
		Joins("INNER JOIN "+tableUnitUserProfileName+" on "+tableUnitUserProfileName+".id = "+tableUnitUserName+".id").
		Where(tableUnitUserName+".unit_id = ?", unitId).
		Where(tableUnitUserName + ".deleted = 0")
		// Where("status <>", base_model.UNIT_USER_PROFILE_CANCLED).
	if id != "" {
		query = query.Where(tableUnitUserName+".id = ?", id)
	} else {
		query = query.Where(tableUnitUserName+".phone = ?", phone)
	}
	err = query.Find(&data).Error

	for i := range data {
		data[i].UnitUser.Id = data[i].Id
		data[i].UnitUserProfile.Id = data[i].Id
	}
	return
}

// 通过id获取组织单位用户
func GetUnitUserById[UnitUserModel itf.UnitUserItf, UnitUserProfileModel itf.UserProfileItf](id string, unitUserModel UnitUserModel, unitUserProfileModel UnitUserProfileModel) (data []unit_dto.UnitUserAllDataDto, err error) {
	tableUnitUserName := unitUserModel.TableName()
	tableUnitUserProfileName := unitUserProfileModel.TableName()

	err = global.GetReadDb().
		Model(unitUserModel).
		Select("*").
		Joins("INNER JOIN "+tableUnitUserProfileName+" on "+tableUnitUserProfileName+".id = "+tableUnitUserName+".id").
		Where(tableUnitUserName+".id = ?", id).
		Where(tableUnitUserName + ".deleted = 0").
		// Where("status <>", base_model.UNIT_USER_PROFILE_CANCLED).
		Find(&data).Error
	return
}

/**
 * 获取用户
 * UnitUserModel: models.PlatUser, models.MchntUser
 * @param userId
 * @param unitId
 * @return
 * @throws
 */
func GetUserOfUnitById[UnitUserModel itf.UnitUserItf](userId string, unitId string) (base_model.UnitUser, error) {
	var unitUserModel UnitUserModel
	tableUnitUserName := unitUserModel.TableName()
	var userData base_model.UnitUser
	if userId == "" {
		return userData, errors.New("userId 不能为空")
	}

	result := global.GetReadDb().
		Model(unitUserModel).
		Select("*, case is_default when 1 then unit_id else '' end AS default_unit_id").
		Where(tableUnitUserName+".user_id = ?", userId).
		Where(tableUnitUserName+".unit_id = ?", unitId).
		Where(tableUnitUserName+".deleted = ?", 0).
		Take(&userData)
	return userData, result.Error
}

// 页面-系统管理-获取用户列表
func GetUserListOfUnitById[
	UnitUserModel itf.UnitUserItf,
	UserProfileModel itf.UserProfileItf,
	UnitDeptModel itf.DeptItf,
	UnitUserDeptModel itf.UserDeptItf,
	UnitRoleModel itf.RoleItf,
	UnitUserRoleModel itf.UserRoleItf,
](
	reqDto page_dto.SystemUserListReqDto,
	unitUserModel UnitUserModel,
	userProfileModel UserProfileModel,
	unitDeptModel UnitDeptModel,
	unitUserDeptModel UnitUserDeptModel,
	unitRoleModel UnitRoleModel,
	unitUserRoleModel UnitUserRoleModel,
) (userData []user_vo.SystemUserListVo, count int64, err error) {
	tableUnitUserName := unitUserModel.TableName()
	tableUnitUserProfileName := userProfileModel.TableName()
	tableUnitDeptName := unitDeptModel.TableName()
	tableUnitUserDeptName := unitUserDeptModel.TableName()
	tableUnitRoleName := unitRoleModel.TableName()
	tableUnitUserRoleName := unitUserRoleModel.TableName()

	tableStruct := struct {
		TableUser        string
		TableUserProfile string
		TableDept        string
		TableUserDept    string
		TableRole        string
		TableUserRole    string
	}{
		TableUser:        tableUnitUserName,
		TableUserProfile: tableUnitUserProfileName,
		TableDept:        tableUnitDeptName,
		TableUserDept:    tableUnitUserDeptName,
		TableRole:        tableUnitRoleName,
		TableUserRole:    tableUnitUserRoleName,
	}

	userData = make([]user_vo.SystemUserListVo, 0)
	if len(reqDto.SelectUnitIds) <= 0 {
		return userData, count, errors.New("GetUserListOfUnitById(): UnitId 不能为空")
	}

	//用户部门子查询
	selectUserDeptTpl := `{{.TableUserDept}}.user_id as unit_user_id, 
		string_agg({{.TableUserDept}}.dept_id, ',') AS dept_ids, 
		string_agg({{.TableDept}}.name, ',') AS dept_names, 
		COUNT(CASE WHEN {{.TableUserDept}}.dept_id IN (?) THEN 1 ELSE null END) AS count_has_id`
	selectUserDeptStr, _ := helper.ParseStringTpl(selectUserDeptTpl, tableStruct)
	subQueryUserDept := global.GetReadDb().
		Model(unitUserDeptModel).
		Select(selectUserDeptStr, reqDto.DeptIds).
		Joins("INNER JOIN "+tableUnitDeptName+" ON "+tableUnitDeptName+".id = "+tableUnitUserDeptName+".dept_id").
		Where(tableUnitDeptName+".deleted = ?", 0).
		Where(tableUnitUserDeptName+".deleted = ?", 0).
		Group(tableUnitUserDeptName + ".user_id")

	//用户角色子查询
	selectUserRoleTpl := `{{.TableUserRole}}.user_id as unit_user_id, 
		string_agg({{.TableUserRole}}.role_id, ',') AS role_ids, 
		string_agg({{.TableRole}}.role_name, ',') AS role_names,
		COUNT(CASE WHEN {{.TableUserRole}}.role_id IN (?) THEN 1 ELSE null END) AS count_has_id`
	selectUserRoleStr, _ := helper.ParseStringTpl(selectUserRoleTpl, tableStruct)
	subQueryUserRole := global.GetReadDb().
		Model(unitUserRoleModel).
		Select(selectUserRoleStr, reqDto.RoleIds).
		Joins("INNER JOIN "+tableUnitRoleName+" ON "+tableUnitRoleName+".id = "+tableUnitUserRoleName+".role_id").
		Where(tableUnitRoleName+".deleted = ?", 0).
		Where(tableUnitUserRoleName+".deleted = ?", 0).
		Group(tableUnitUserRoleName + ".user_id")

	selectStr := "*"
	query := global.GetReadDb().
		Model(unitUserModel).
		Joins("INNER JOIN "+tableUnitUserProfileName+" ON "+tableUnitUserProfileName+".id = "+tableUnitUserName+".id").
		Joins("LEFT JOIN (?) t ON t.unit_user_id = "+tableUnitUserName+".id", subQueryUserDept).
		Joins("LEFT JOIN (?) t2 ON t2.unit_user_id = "+tableUnitUserName+".id", subQueryUserRole).
		Where(tableUnitUserName+".unit_id = ?", reqDto.UnitId).
		Where(tableUnitUserName+".deleted = ?", 0)

	if len(reqDto.DeptIds) > 0 {
		query = query.Where("t.count_has_id>0")
	}
	if len(reqDto.RoleIds) > 0 {
		query = query.Where("t2.count_has_id>0")
	}

	err = query.Select(tableUnitUserName + ".id").Count(&count).Error
	if err != nil {
		return userData, count, err
	}
	if count == 0 {
		return userData, count, nil
	}

	result := query.
		Select(selectStr).
		Limit(reqDto.PageSize).
		Offset(reqDto.Offset).
		Find(&userData)

	return userData, count, result.Error
}

/**
 * 获取用户默认组织单位
 * UnitModel: models.PlatUnit, models.MchntUnit
 * UnitUserModel: models.PlatUser, models.MchntUser
 * @param userId
 * @return
 * @throws
 */
func GetUserDefaultUnit[UnitModel itf.UnitItf, UnitUserModel itf.UnitUserItf](userId string) (unit_vo.UnitListVo, error) {
	var unitModel UnitModel
	var unitUserModel UnitUserModel
	tableUnitName := unitModel.TableName()
	tableUnitUserName := unitUserModel.TableName()
	userData := unit_vo.UnitListVo{}
	if userId == "" {
		return userData, errors.New("userId 不能为空")
	}

	tableStruct := struct {
		TableUnit     string
		TableUserUnit string
	}{
		TableUnit:     tableUnitName,
		TableUserUnit: tableUnitUserName,
	}
	joinUserUnitStr, err1 := helper.ParseStringTpl(`inner join {{.TableUserUnit}} on {{.TableUserUnit}}.unit_id = {{.TableUnit}}.id`, tableStruct)
	selectStr, err2 := helper.ParseStringTpl(`{{.TableUnit}}.*, case {{.TableUserUnit}}.is_default when 1 then {{.TableUserUnit}}.unit_id else '' end AS default_unit_id, case {{.TableUserUnit}}.is_default when 1 then {{.TableUserUnit}}.id else '' end AS default_unit_user_id`, tableStruct)
	if err1 != nil {
		return userData, err1
	}
	if err2 != nil {
		return userData, err2
	}

	result := global.GetReadDb().
		Model(unitModel).
		Select(selectStr).
		Joins(joinUserUnitStr).
		Where(tableUnitUserName+".user_id = ?", userId).
		Where(tableUnitUserName+".is_default = ?", 1).
		Where(tableUnitUserName+".deleted = ?", 0).
		Take(&userData)
	return userData, result.Error
}

/**
 * 更新用户默认组织单位id
 * UnitUserModel: models.PlatUser, models.MchntUser
 * @param userId
 * @param unitId
 * @return
 * @throws
 */
func UpdateUserDefaultUnit[UnitUserModel itf.UnitUserItf](userId string, unitId string) error {
	var unitUserModel UnitUserModel
	_, err := GetUserOfUnitById[UnitUserModel](userId, unitId)
	if err != nil {
		return err
	}

	updateData := struct {
		IsDefault int
	}{
		IsDefault: 0,
	}
	err = global.GetWriteDb().Transaction(func(tx *gorm.DB) error {
		result := tx.Model(unitUserModel).Select("is_default").Where("user_id = ?", userId).Updates(updateData)
		if result.Error != nil {
			return result.Error
		}

		updateData.IsDefault = 1
		result = tx.Model(unitUserModel).Where("user_id = ?", userId).Where("unit_id = ?", unitId).Updates(updateData)
		if result.Error != nil {
			return result.Error
		}
		return nil
	})
	return err
}

/**
 * 新增默认组织单位-新增其默认用户
 * UnitUserModel: models.PlatUser, models.MchntUser
 * @param userId
 * @param unitId
 * @param isDefault
 * @param isAdmin
 * @return
 * @throws
 */
func InsertUnitUserForCreateUnit[UnitUserModel itf.UnitUserItf](tx *gorm.DB, userId string, unitId string, isDefault int, unitUserModel UnitUserModel) (unitUserTableUuid string, err error) {
	tableName := unitUserModel.TableName()
	fmt.Println(tableName)
	unitUserTableUuid, _ = helper.GetUuid()
	userInfo := &models.User{}
	global.GetReadDb().Model(&models.User{}).Where("id = ?", userId).Take(&userInfo)

	var insertUnitUserData = base_model.UnitUser{
		Id:        unitUserTableUuid,
		UserId:    userId,
		UnitId:    unitId,
		IsDefault: isDefault,
		Deleted:   0,
		Name:      userInfo.Name,
		Phone:     userInfo.Phone,
	}
	// user 组织单位用户
	err = tx.Model(unitUserModel).
		Create(&insertUnitUserData).Error
	if err != nil {
		return
	}
	return
}
func UpsertUnitUser[UnitUserModel itf.UnitUserItf](tx *gorm.DB, saveData base_model.UnitUser, unitUserModel UnitUserModel) (unitUserTableUuid string, err error) {
	tableName := unitUserModel.TableName()
	fmt.Println(tableName)
	if saveData.Id == "" {
		unitUserTableUuid, _ = helper.GetUuid()
		saveData.Id = unitUserTableUuid
	} else {
		unitUserTableUuid = saveData.Id
	}

	// user 组织单位用户
	err = tx.Model(unitUserModel).
		Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "id"}},
			DoUpdates: clause.AssignmentColumns([]string{"unit_id" /*"phone",*/, "name"}), // 手机号
		}).
		Create(&saveData).Error
	if err != nil {
		return
	}
	return
}

// 验证组织结构
func CheckOrgStructure[UnitUserModel itf.UnitUserItf, UnitDeptModel itf.DeptItf, UnitRoleModel itf.RoleItf](result interface{}, unitUserId string, unitId string, deptId string, roleId []string, unitUserModel UnitUserModel, unitDeptModel UnitDeptModel, unitRoleModel UnitRoleModel) (err error) {
	tableUnitUserName := unitUserModel.TableName()
	tableUnitDeptName := unitDeptModel.TableName()
	tableUnitRoleName := unitRoleModel.TableName()

	tableStruct := struct {
		TableUnitUser string
		TableUnitDept string
		TableUnitRole string
	}{
		TableUnitUser: tableUnitUserName,
		TableUnitDept: tableUnitDeptName,
		TableUnitRole: tableUnitRoleName,
	}

	sqlTpl := `
		SELECT 
			SUM(existUnit) AS exist_unit,
			SUM(existDept) AS exist_dept,
			SUM(existRole) AS exist_role
		FROM (
			(SELECT count(1) existUnit, 0 as existDept, 0 as existRole FROM {{.TableUnitUser}} WHERE {{.TableUnitUser}}.id = ? AND {{.TableUnitUser}}.unit_id = ? AND {{.TableUnitUser}}.deleted = 0) UNION ALL
			(SELECT 0 existUnit, count(1) as existDept, 0 as existRole FROM {{.TableUnitDept}} WHERE {{.TableUnitDept}}.id = ? AND {{.TableUnitDept}}.unit_id = ? AND {{.TableUnitDept}}.deleted = 0) UNION ALL
			(SELECT 0 as existUnit, 0 as existDept, count(1) existRole FROM {{.TableUnitRole}} WHERE {{.TableUnitRole}}.id IN ? AND {{.TableUnitRole}}.unit_id = ? AND {{.TableUnitRole}}.deleted = 0)
		) AS t 
		
	`
	sql, err := helper.ParseStringTpl(sqlTpl, tableStruct)
	if err != nil {
		return err
	}

	tmpsql := global.GetReadDb().Raw(sql, unitUserId, unitId, deptId, unitId, roleId, unitId).Statement.SQL.String()
	fmt.Println(tmpsql)
	return global.GetReadDb().Raw(sql, unitUserId, unitId, deptId, unitId, roleId, unitId).Scan(&result).Error
}

// 删除组织单位用户

func DelUnitUser[UnitUserModel itf.UnitUserItf, UnitUserProfileModel itf.UserProfileItf](ids []string, unitUserModel UnitUserModel, unitUserProfileModel UnitUserProfileModel) (err error) {
	upUnitUserData := map[string]interface{}{
		"deleted": 1,
	}
	upUnitUserProfileData := map[string]interface{}{
		"deleted":    1,
		"deleted_at": helper.GetTimestamp(),
	}
	err = global.WriteDb.Transaction(func(tx *gorm.DB) (err error) {
		err = tx.Model(unitUserModel).Where("id IN ?", ids).Updates(upUnitUserData).Error
		if err != nil {
			return
		}
		err = tx.Model(unitUserProfileModel).Where("id IN ?", ids).Updates(upUnitUserProfileData).Error
		return
	})
	return
}
