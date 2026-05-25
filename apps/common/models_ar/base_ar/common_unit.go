package base_ar

import (
	"WenBeego/apps/common/dto_vo/page_dto"
	"WenBeego/apps/common/dto_vo/unit_dto"
	"WenBeego/apps/common/dto_vo/unit_vo"
	"WenBeego/apps/common/global"
	"WenBeego/apps/common/helper"
	"WenBeego/apps/common/models/base_model"
	"WenBeego/apps/common/models/itf"
	"errors"
	"strings"

	"gorm.io/gorm"
)

/**
 * 获取所有组织单位
 */
func GetAllUnit[UnitModel itf.UnitItf](unitModel UnitModel, fields string) ([]base_model.Unit, error) {
	if fields == "" {
		fields = "*"
	}
	listData := []base_model.Unit{}
	result := global.GetReadDb().
		Model(unitModel).
		Select(fields).
		Where(unitModel.TableName() + ".deleted = 0").
		Find(&listData)
	return listData, result.Error
}

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

	selectStr, err := helper.ParseStringTpl(`{{.TableUnit}}.id,{{.TableUnit}}.pid,{{.TableUnit}}.name,{{.TableUnit}}.logo,{{.TableUnit}}.status`, tableStruct)
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
		Where(tableUnitName + ".deleted = 0").
		Order(tableUnitName + ".created_by," + tableUnitName + ".sort").
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
		Where(tableUnitName + ".deleted = 0").
		Take(&unitModel)
	return unitModel, result.Error
}

// 获取内部组织列表
func GetUnitListByUserId[UnitModel itf.UnitItf, UnitUserModel itf.UnitUserItf](unitDto page_dto.SystemUnitListReqDto, unitModel UnitModel, unitUserModel UnitUserModel) (listData []unit_vo.UnitListVo, count int64, err error) {
	tableUnitName := unitModel.TableName()
	tableUnitUserName := unitUserModel.TableName()

	unitDto.Name = strings.TrimSpace(unitDto.Name)
	unitDto.Code = strings.TrimSpace(unitDto.Code)
	listData = make([]unit_vo.UnitListVo, 0)

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
		Where(tableUnitUserName + ".deleted = 0").
		Where(tableUnitName + ".deleted = 0")

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
	if count == 0 {
		return nil, 0, nil
	}

	err = query.Select(selectStr).Order(tableUnitName + ".created_by," + tableUnitName + ".sort").Find(&listData).Error
	return
}
func GetUnitListForAdminPlat[UnitModel itf.UnitItf, UnitUserModel itf.UnitUserItf](unitDto page_dto.SystemUnitListReqDto, unitModel UnitModel, unitUserModel UnitUserModel) (listData []base_model.Unit, count int64, err error) {
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
		// Where(tableUnitUserName+".user_id = ?", unitDto.UserId).
		Where(tableUnitUserName + ".deleted = 0").
		Where(tableUnitName + ".deleted = 0")

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
	if count == 0 {
		return nil, 0, nil
	}

	err = query.Select(selectStr).Order(tableUnitName + ".created_by," + tableUnitName + ".sort").Find(&listData).Error
	return
}

func SaveUnit[UnitModel itf.UnitItf](tx *gorm.DB, unitDto unit_dto.UnitDto, unitModel UnitModel) (id string, err error) {
	if unitDto.Id == "" {
		unitDto.IsOfficial = false
		unitDto.Id, err = helper.GetUuid()
		if err != nil {
			return
		}

		err = tx.Model(unitModel).
			Create(&unitDto).Error
		if err != nil {
			return
		}
	} else {
		err = tx.Model(unitModel).
			Where("id = ?", unitDto.Id).
			Select("*").
			Omit("is_official", "deleted", "created_at", "deleted_at").
			Updates(&unitDto).Error
	}
	if err != nil {
		return
	}
	return unitDto.Id, nil
}

func DeleteUnit[UnitModel itf.UnitItf](id string, unitModel UnitModel) error {
	deletedTime := helper.GetTimestamp()

	updateData := struct {
		Deleted   int
		UpdatedAt int64
		DeletedAt *int64
	}{
		Deleted:   1,
		UpdatedAt: deletedTime,
		DeletedAt: &deletedTime,
	}

	return global.GetWriteDb().
		Model(unitModel).
		Where("id = ?", id).
		Updates(updateData).Error
}

// 删除组织单位
func DelUnit[UnitModel itf.UnitItf](id string, unitUserId string, unitModel UnitModel) error {
	if id == "" {
		return errors.New("DelUnit：参数id不能为空")
	}
	timeInt := helper.GetTimestamp()
	updateData := struct {
		Deleted   int
		UpdatedAt int64
		DeletedAt *int64
		UpdatedBy string
		DeletedBy string
	}{
		Deleted:   1,
		UpdatedAt: timeInt,
		DeletedAt: &timeInt,
		UpdatedBy: unitUserId,
		DeletedBy: unitUserId,
	}
	return global.GetWriteDb().
		Model(unitModel).
		Where("id = ?", id).
		Updates(updateData).Error
}
