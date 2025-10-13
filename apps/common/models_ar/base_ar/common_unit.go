package base_ar

import (
	"WenBeego/apps/common/dto/page_dto"
	"WenBeego/apps/common/global"
	"WenBeego/apps/common/helper"
	"WenBeego/apps/common/models/base_model"
	"WenBeego/apps/common/models/itf"
	"strings"
)

/**
 * 获取用户组织列表
 * @param {string} userId 用户ID
 * @param {models.ModelInterface} unitModel models.Plat, models.Mchnt
 * @param {models.ModelInterface} unitUserModel models.PlatUser, models.MchntUser
 * @return {[]models.ModelInterface} []models.Plat, []models.Mchnt
 * @return {error}
 */
func GetUserUnitList[UnitModel itf.UnitItf, UnitUserModel itf.UnitUserItf](userId string, unitModel UnitModel, unitUserModel UnitUserModel) ([]UnitModel, error) {
	tableUnitName := unitModel.TableName()
	tableUnitUserName := unitUserModel.TableName()

	tableStruct := struct {
		TableUnit     string
		TableUserUnit string
	}{
		TableUnit:     tableUnitName,
		TableUserUnit: tableUnitUserName,
	}

	selectStr, err := helper.ParseStringTpl(`{{.TableUnit}}.*`, tableStruct)
	joinStr, err2 := helper.ParseStringTpl(`inner join {{.TableUserUnit}} on {{.TableUserUnit}}.unit_id = {{.TableUnit}}.id`, tableStruct)
	if err != nil {
		return nil, err
	}
	if err2 != nil {
		return nil, err2
	}

	listData := []UnitModel{}
	result := global.GetReadDb().
		Model(unitModel).
		Select(selectStr).
		Joins(joinStr).
		Where(tableUnitUserName+".user_id = ?", userId).
		Where(tableUnitUserName + ".deleted = 0").
		Find(&listData)
	return listData, result.Error
}

// 查找是否存在用户组织单位
func GetUserUnitById[UnitModel itf.UnitItf, UnitUserModel itf.UnitUserItf](userId string, unitId string, unitModel UnitModel, unitUserModel UnitUserModel) (UnitModel, error) {
	tableUnitName := unitModel.TableName()
	tableUnitUserName := unitUserModel.TableName()

	tableStruct := struct {
		TableUnit     string
		TableUserUnit string
	}{
		TableUnit:     tableUnitName,
		TableUserUnit: tableUnitUserName,
	}

	selectStr, err := helper.ParseStringTpl(`{{.TableUnit}}.*`, tableStruct)
	joinStr, err2 := helper.ParseStringTpl(`inner join {{.TableUserUnit}} on {{.TableUserUnit}}.unit_id = {{.TableUnit}}.id`, tableStruct)
	if err != nil {
		return unitModel, err
	}
	if err2 != nil {
		return unitModel, err2
	}

	result := global.GetReadDb().
		Model(unitModel).
		Select(selectStr).
		Joins(joinStr).
		Where(tableUnitUserName+".user_id = ?", userId).
		Where(tableUnitUserName+".unit_id = ?", unitId).
		Where(tableUnitUserName + ".deleted = 0").
		Take(&unitModel)
	return unitModel, result.Error
}

// 获取内部组织列表
func GetUnitListById[UnitModel itf.UnitItf, UnitUserModel itf.UnitUserItf](unitDto page_dto.SystemUnitListReqDto, unitModel UnitModel, unitUserModel UnitUserModel) (listData []base_model.Unit, count int64, err error) {
	tableUnitName := unitModel.TableName()
	tableUnitUserName := unitUserModel.TableName()

	unitDto.Name = strings.TrimSpace(unitDto.Name)
	unitDto.Code = strings.TrimSpace(unitDto.Code)

	tableStruct := struct {
		TableUnit     string
		TableUserUnit string
	}{
		TableUnit:     tableUnitName,
		TableUserUnit: tableUnitUserName,
	}

	selectStr, err := helper.ParseStringTpl(`{{.TableUnit}}.*`, tableStruct)
	joinStr, err2 := helper.ParseStringTpl(`inner join {{.TableUserUnit}} on {{.TableUserUnit}}.unit_id = {{.TableUnit}}.id`, tableStruct)
	if err != nil {
		return nil, 0, err
	}
	if err2 != nil {
		return nil, 0, err2
	}

	query := global.GetReadDb().
		Model(unitModel).
		Joins(joinStr).
		Where(tableUnitUserName+".user_id = ?", unitDto.UserId).
		Where(tableUnitUserName + ".deleted = 0")

	if unitDto.ParentUnitId != "" {
		query = query.Where(tableUnitName+".pid = ?", unitDto.ParentUnitId)
	} else {
		query = query.Where(tableUnitName + ".pid = ''")
	}
	if unitDto.Name != "" {
		query = query.Where(tableUnitName+".name like ?", "%"+unitDto.Name+"%")
	}
	if unitDto.Code != "" {
		query = query.Where(tableUnitName+".code like ?", "%"+unitDto.Code+"%")
	}
	if unitDto.Status != -1 {
		query = query.Where(tableUnitName+".status = ?", unitDto.Status)
	}

	err = query.Select(tableUnitName + ".id").Count(&count).Error
	if err != nil {
		return nil, 0, err
	}

	err = query.Select(selectStr).Find(&listData).Error
	return
}
