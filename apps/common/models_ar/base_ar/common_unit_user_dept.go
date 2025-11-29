package base_ar

import (
	"WenBeego/apps/common/dto/user_dto"
	"WenBeego/apps/common/helper"
	"WenBeego/apps/common/models/itf"
	"errors"

	"gorm.io/gorm"
)

func SaveUnitUserDept[UserDeptModel itf.UserDeptItf](tx *gorm.DB, userDeptData user_dto.UnitUserDeptDto, userDeptModel UserDeptModel) (err error) {
	if userDeptData.Id == "" {
		userDeptData.Id, err = helper.GetUuid()
		if err != nil {
			return
		}
	}
	if userDeptData.UserId == "" || userDeptData.DeptId == "" {
		return errors.New("新增用户角色，用户id和角色id不能为空")
	}
	err = tx.Where("user_id = ? AND dept_id = ? ", userDeptData.UserId, userDeptData.DeptId).Delete(&userDeptModel).Error
	if err != nil {
		return
	}
	err = tx.Model(userDeptModel).
		Create(&userDeptData).Error
	return
}
