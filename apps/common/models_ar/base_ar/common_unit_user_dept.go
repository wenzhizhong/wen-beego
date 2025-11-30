package base_ar

import (
	"WenBeego/apps/common/dto/user_dto"
	"WenBeego/apps/common/helper"
	"WenBeego/apps/common/models/itf"
	"errors"
	"fmt"

	"gorm.io/gorm"
)

func SaveUnitUserDept[DeptModel itf.DeptItf, UserDeptModel itf.UserDeptItf](tx *gorm.DB, unit_id string, userDeptData user_dto.UnitUserDeptDto, deptModel DeptModel, userDeptModel UserDeptModel) (err error) {
	if userDeptData.Id == "" {
		userDeptData.Id, err = helper.GetUuid()
		if err != nil {
			return
		}
	}
	if userDeptData.UserId == "" || userDeptData.DeptId == "" {
		return errors.New("新增用户角色，用户id和角色id不能为空")
	}
	tableDeptName := deptModel.TableName()
	// tableUserDeptName := userDeptModel.TableName()
	err = tx.
		Where(fmt.Sprintf("user_id = ? AND dept_id IN (SELECT id FROM %s WHERE unit_id = ?)", tableDeptName), userDeptData.UserId, unit_id).
		Delete(&userDeptModel).Error
	if err != nil {
		return
	}
	err = tx.Model(userDeptModel).
		Create(&userDeptData).Error
	return
}
