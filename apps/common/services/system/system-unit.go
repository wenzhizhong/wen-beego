package system

import (
	"WenBeego/apps/common/dto"
	"WenBeego/apps/common/dto/page_dto"
	"WenBeego/apps/common/dto/unit_dto"
	"WenBeego/apps/common/global"
	"WenBeego/apps/common/helper"
	"WenBeego/apps/common/middleware/business_store"
	"WenBeego/apps/common/models"
	"WenBeego/apps/common/models/base_model"
	"WenBeego/apps/common/models/itf"
	"WenBeego/apps/common/models_ar/base_ar"
	"errors"

	"gorm.io/gorm"
)

type Unit struct {
}

// 系统管理-获取内部组织列表
func (s *Unit) GetUnitList(unitDto page_dto.SystemUnitListReqDto) (resultDto dto.RespDataListDto, err error) {
	data := make([]base_model.Unit, 0)
	var count int64 = 0

	switch unitDto.ModuleName {
	case "admin_plat":
		data, count, err = base_ar.GetUnitListByUserId(unitDto, &models.Plat{}, &models.PlatUser{})
	case "mchnt_plat":
		data, count, err = base_ar.GetUnitListByUserId(unitDto, &models.Mchnt{}, &models.MchntUser{})
	default:
		err = errors.New("GetUnitList:模块名称错误")
	}
	if err != nil {
		return
	}
	for k, v := range data {
		tmpLogo, err1 := helper.LocalFileSign(unitDto.Host, v.Logo)
		tmpPath, err2 := helper.LocalFileSign(unitDto.Host, v.License)
		if err1 != nil || err2 != nil {
			errStr := helper.Ternary(err1 != nil, err1.Error(), err2.Error())
			global.Log.Error(errStr)
			continue
		}
		data[k].LogoLink = tmpLogo
		data[k].LicenseLink = tmpPath
	}
	resultDto, err = helper.GetRespDataListDto(unitDto.PageSize, unitDto.CurrentPage, count, data)
	return
}

// 系统管理-保存用户
func (s *Unit) Save(baseParamDto dto.BaseParamDto, unitDto unit_dto.UnitDto) (result map[string]string, err error) {
	newUnitId := ""
	isAdd := unitDto.Id == ""
	if isAdd {
		unitDto.CreatedBy = baseParamDto.UnitUserId
	} else {
		unitDto.UpdatedBy = baseParamDto.UnitUserId
	}

	switch baseParamDto.ModuleName {
	case "admin_plat":
		newUnitId, err = doSave(isAdd, baseParamDto, unitDto, &models.Plat{}, &models.PlatUser{}, &models.PlatUserProfile{}, &models.PlatRole{}, &models.PlatUserRole{}, &models.PlatMenu{}, &models.PlatRoleClassify{})
	case "mchnt_plat":
		newUnitId, err = doSave(isAdd, baseParamDto, unitDto, &models.Mchnt{}, &models.MchntUser{}, &models.MchntUserProfile{}, &models.MchntRole{}, &models.MchntUserRole{}, &models.MchntMenu{}, &models.MchntRoleClassify{})
	default:
		err = errors.New("unit Save：模块名称错误")
	}
	if err != nil {
		global.Log.Error(err.Error())
		return
	}
	// 清空用户权限认证缓存
	business_store.ClearAumid()

	result = make(map[string]string)
	result["id"] = newUnitId
	return
}
func (s *Unit) Del(baseParamDto dto.BaseParamDto, unitDto unit_dto.UnitDto) (err error) {
	switch baseParamDto.ModuleName {
	case "admin_plat":
		err = base_ar.DelUnit(unitDto.Id, baseParamDto.UnitUserId, &models.Plat{})
	case "mchnt_plat":
		err = base_ar.DelUnit(unitDto.Id, baseParamDto.UnitUserId, &models.Mchnt{})
	default:
		err = errors.New("Unit Del:模块名称错误")
	}
	return
}
func doSave[
	UnitModel itf.UnitItf,
	UnitUserModel itf.UnitUserItf,
	UnitUserProfileModel itf.UserProfileItf,
	RoleModel itf.RoleItf,
	UnitUserRoleModel itf.UserRoleItf,
	UnitMenuModel itf.MenuItf,
	RoleClassifyModel itf.RoleClassifyItf,
](
	isAdd bool,
	baseParamDto dto.BaseParamDto,
	unitDto unit_dto.UnitDto,
	unitModel UnitModel,
	unitUserModel UnitUserModel,
	unitUserProfileModel UnitUserProfileModel,
	roleModel RoleModel,
	unitUserRoleModel UnitUserRoleModel,
	unitMenuModel UnitMenuModel,
	roleClassifyModel RoleClassifyModel,
) (newUnitId string, err error) {

	err = global.GetWriteDb().Transaction(func(tx *gorm.DB) error {
		newUnitId, err = base_ar.SaveUnit(tx, unitDto, unitModel)
		if err != nil {
			return err
		}

		if isAdd {
			unitUserTableUuid := ""
			unitUserTableUuid, err = base_ar.InsertUnitUserForCreateUnit(tx, baseParamDto.UserId, newUnitId, 0, unitUserModel)
			if err != nil {
				return err
			}

			userProfile := getUserProfileObj(unitUserTableUuid)
			err = base_ar.InsertUnitUserProfileForCreateUnit(tx, userProfile, unitUserProfileModel)
			if err != nil {
				return err
			}

			unitRole, err := getUnitRoleObj(newUnitId, baseParamDto.UserId)
			if err != nil {
				return err
			}
			err = base_ar.InsertUnitRole[RoleModel](tx, unitRole)
			if err != nil {
				return err
			}

			unitRoleClassify, err := getUnitRoleClassifyObj(newUnitId, unitRole.Id)
			if err != nil {
				return err
			}
			err = base_ar.InsertUserRoleClassifies[RoleClassifyModel](tx, unitRoleClassify)
			if err != nil {
				return err
			}

			unitUserRole, err := getUnitUserRoleObj(baseParamDto.UserId, unitRole.Id)
			if err != nil {
				return err
			}
			err = base_ar.InsertUnitUserRole[UnitUserRoleModel](tx, unitUserRole)
			if err != nil {
				return err
			}

			err = base_ar.CloneMenu[UnitMenuModel](tx, baseParamDto.UnitId, newUnitId)
			if err != nil {
				return err
			}

		}
		return nil
	})
	return
}

func getUserProfileObj(unitUserId string) base_model.UnitUserProfile {

	return base_model.UnitUserProfile{
		Id:             unitUserId,
		Avatar:         "",
		CardType:       base_model.UNIT_CARD_TYPE_5,
		CardNum:        "",
		CardImages:     "",
		Gender:         0,
		BirthDate:      nil,
		Constellation:  "",
		Occupation:     "",
		Company:        "",
		EmergencyName:  "",
		EmergencyTel:   "",
		Address:        "",
		Email:          "",
		Source:         base_model.UNIT_USER_SOURCE_OTHER,
		ValidDateBegin: nil,
		ValidDateEnd:   nil,
		Schooling:      "",
		DegreeNumber:   "",
		Professional:   "",
		Status:         base_model.UNIT_USER_PROFILE_NORMAL,
		// CreatedAt:    time.Now().Unix(),
		// UpdatedAt:    0,
		Deleted:   0,
		DeletedAt: nil,
	}
}

func getUnitRoleObj(newUnitId string, userId string) (data base_model.UnitRole, err error) {
	uuid, err := helper.GetUuid()
	if err != nil {
		return data, err
	}
	data = base_model.UnitRole{
		Id:        uuid,
		UnitId:    newUnitId,
		RoleName:  "超级管理员",
		RoleSort:  1,
		Status:    1,
		Deleted:   0,
		CreatedBy: userId,
		// CreatedAt: ,
		UpdatedBy: userId,
		// UpdatedAt: ,
		Remark: "",
	}
	return data, nil
}

func getUnitRoleClassifyObj(newUnitId string, roleId string) (data base_model.UnitRoleClassify, err error) {
	uuid, err := helper.GetUuid()
	data = base_model.UnitRoleClassify{
		Id:      uuid,
		RoleId:  roleId,
		UnitId:  newUnitId,
		Name:    "admin",
		Deleted: 0,
	}
	return data, err
}
func getUnitUserRoleObj(userId string, roleId string) (data base_model.UnitUserRole, err error) {
	uuid, err := helper.GetUuid()
	data = base_model.UnitUserRole{
		Id:      uuid,
		UserId:  userId,
		RoleId:  roleId,
		Deleted: 0,
	}
	return data, err
}
