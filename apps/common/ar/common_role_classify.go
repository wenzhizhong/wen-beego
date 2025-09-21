package ar

import (
	"WenBeego/apps/common/global"
	"WenBeego/apps/common/helper"
	"WenBeego/apps/common/models"
	"errors"
	"fmt"
)

func GetUserRoleClassifies[UnitModel models.ModelInterface, RoleModel models.ModelInterface, RoleClassifyModel models.ModelInterface, UserRoleModel models.ModelInterface](unitId string, userId string, unitModel UnitModel, roleModel RoleModel, roleClassifyModel RoleClassifyModel, userRoleModel UserRoleModel) (dataList []map[string]interface{}, err error) {
	if unitId == "" || userId == "" {
		str := fmt.Sprintf("GetUserRoleClassifies():获取菜单权限必填参数, unit_id:%s, classifyName:%s", unitId, userId)
		global.Log.Error(str)
		return dataList, errors.New(str)
	}

	tableUnit := unitModel.TableName()
	tableUserRole := userRoleModel.TableName()
	tableRole := roleModel.TableName()
	tableRoleClassify := roleClassifyModel.TableName()
	tableStruct := struct {
		TableUnit         string
		TableRole         string
		TableRoleClassify string
		TableUserRole     string
	}{
		TableUnit:         tableUnit,
		TableRole:         tableRole,
		TableRoleClassify: tableRoleClassify,
		TableUserRole:     tableUserRole,
	}

	selectStr, err := helper.ParseStringTpl(`{{.TableRoleClassify}}.*`, tableStruct)
	joinUserRoleStr, err3 := helper.ParseStringTpl(`inner join {{.TableUserRole}} on {{.TableUserRole}}.role_id = {{.TableRoleClassify}}.role_id`, tableStruct)
	joinRoleStr, err4 := helper.ParseStringTpl(`inner join {{.TableRole}} on {{.TableRole}}.id = {{.TableRoleClassify}}.role_id`, tableStruct)
	joinUnitStr, err5 := helper.ParseStringTpl(`inner join {{.TableUnit}} on {{.TableUnit}}.id = {{.TableUserRole}}.unit_id`, tableStruct)
	if err != nil {
		return dataList, err
	}
	if err3 != nil {
		return dataList, err3
	}
	if err4 != nil {
		return dataList, err4
	}
	if err5 != nil {
		return dataList, err5
	}

	tmpError := global.GetReadDb().
		Model(&roleClassifyModel).
		Select(selectStr).
		Joins(joinRoleStr).
		Joins(joinUserRoleStr).
		Joins(joinUnitStr).
		Where(tableRoleClassify+".unit_id = ?", unitId).
		Where(tableRoleClassify+".deleted = ?", 0).
		Where(tableRole+".deleted = ?", 0).
		Where(tableRole+".status = ?", 1).
		Where(tableUserRole+".user_id = ?", userId).
		Where(tableUserRole+".unit_id = ?", unitId).
		Where(tableUserRole+".deleted = ?", 0).
		Where(tableUnit+".deleted = ?", 0).
		Where(tableUnit+".status = ?", 1).
		Scan(&dataList).Error
	if tmpError != nil && !helper.DbNotFound(tmpError) {
		err = tmpError
		return
	}
	return
}
