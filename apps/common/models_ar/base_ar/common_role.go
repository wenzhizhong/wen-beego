package base_ar

import (
	"WenBeego/apps/common/dto/page_dto"
	"WenBeego/apps/common/dto/role_dto"
	"WenBeego/apps/common/global"
	"WenBeego/apps/common/models/base_model"
	"WenBeego/apps/common/models/itf"
	"errors"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// 页面-保存部门
func SaveUnitRole[UnitRoleModel itf.RoleItf](tx *gorm.DB, unitRoleDto role_dto.UnitRoleDto, unitRoleModel UnitRoleModel) (id string, err error) {
	err = tx.Model(unitRoleModel).
		Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "id"}},
			UpdateAll: true,
		}).
		Create(&unitRoleDto).Error
	return unitRoleDto.Id, err
}

// 页面-删除组织架构
func DelUnitRole[UnitRoleModel itf.RoleItf](unitRoleData base_model.UnitRole, unitRoleModel UnitRoleModel) error {
	if unitRoleData.Id == "" {
		return errors.New("DelUnitRole: 参数id不能为空")
	}
	return global.GetWriteDb().
		Model(unitRoleModel).
		Where("id = ?", unitRoleData.Id).
		Updates(unitRoleData).Error
}

/**
 * 新增组织单位-初始化新增组织单位角色配置
 */
func InsertUnitRole[RoleModel itf.RoleItf](tx *gorm.DB, roleModel base_model.UnitRole) (err error) {
	if roleModel.Id == "" {
		return errors.New("新增角色，角色id不能为空")
	}
	var tmpUnitRole RoleModel
	err = global.GetReadDb().
		Model(tmpUnitRole).
		Where("id = ?", roleModel.Id).
		Take(&tmpUnitRole).Error
	if err == nil && tmpUnitRole.GetId() != "" {
		return nil
	}

	err = tx.Model(tmpUnitRole).
		Create(&roleModel).Error
	return
}

// 页面-获取组织架构部门列表
func GetUnitRoleList[UnitModel itf.UnitItf, UnitRoleModel itf.RoleItf, UnitRoleClassifyModel itf.RoleClassifyItf](RoleDto page_dto.SystemRoleListReqDto, unitModel UnitModel, unitRoleModel UnitRoleModel, unitRoleClassifyModel UnitRoleClassifyModel) (dataList []base_model.UnitRole, count int64, err error) {
	tableUnitName := unitModel.TableName()
	tableRoleName := unitRoleModel.TableName()
	tableRoleClassifyName := unitRoleClassifyModel.TableName()

	query := global.GetReadDb().
		Model(unitRoleModel).
		Joins("INNER JOIN " + tableUnitName + " ON " + tableUnitName + ".id = " + tableRoleName + ".unit_id").
		Joins("INNER JOIN " + tableRoleClassifyName + " ON " + tableRoleClassifyName + ".role_id = " + tableRoleName + ".id").
		Where(tableRoleName + ".deleted=0").
		Where(tableUnitName + ".deleted=0")

	if RoleDto.RoleName != "" {
		query = query.Where(tableRoleName+".role_name like ?", "%"+RoleDto.RoleName+"%")
	}
	if len(RoleDto.SelectUnitIds) > 0 {
		query = query.Where(tableRoleName+".unit_id in ?", RoleDto.SelectUnitIds)
	}
	if RoleDto.RoleClassifyName != "" {
		query = query.Where(tableRoleClassifyName+".name like ?", "%"+RoleDto.RoleClassifyName+"%")
	}
	if RoleDto.Status != -1 {
		query = query.Where(tableRoleName+".status = ?", RoleDto.Status)
	}

	err = query.Select(tableRoleName + ".id").Count(&count).Error
	if err != nil {
		return make([]base_model.UnitRole, 0), 0, nil
	}
	if count == 0 {
		return make([]base_model.UnitRole, 0), 0, nil
	}
	err = query.Select(tableRoleName + ".*," + tableUnitName + ".name as unit_name," + tableRoleClassifyName + ".name as role_classify_name").
		Order(tableRoleName + ".role_sort").
		Limit(RoleDto.PageSize).
		Offset(RoleDto.Offset).
		Find(&dataList).Error
	return
}

// 获取组织架构树
func GetUnitRoleTree[UnitRoleModel itf.RoleItf](selectUnitIds []string, unitRoleModel UnitRoleModel) (data []base_model.UnitRole, err error) {
	data = make([]base_model.UnitRole, 0)
	if len(selectUnitIds) == 0 {
		return data, nil
	}
	// tableRoleName := unitRoleModel.TableName()
	err = global.GetReadDb().
		Model(unitRoleModel).
		Select("id,role_name,status").
		Where("unit_id in (?)", selectUnitIds).
		Where("deleted = 0").
		Order("role_sort").
		Find(&data).Error
	return
}
