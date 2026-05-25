package system

import (
	dto "WenBeego/apps/common/dto_vo"
	"WenBeego/apps/common/dto_vo/menu_dto"
	"WenBeego/apps/common/dto_vo/page_dto"
	"WenBeego/apps/common/global"
	"WenBeego/apps/common/helper"
	"WenBeego/apps/common/models"
	"WenBeego/apps/common/models/base_model"
	CommonSystem "WenBeego/apps/common/models_ar"
	"WenBeego/apps/common/models_ar/base_ar"
	"errors"

	"gorm.io/gorm"
)

type MenuService struct {
	commonPlatMenuAr  CommonSystem.PlatMenuAr
	commonMchntMenuAr CommonSystem.MchntMenuAr
}

/**
 * 获取菜单列表
 * @param pageDto
 * @param platformType "admin_plat", "admin_mchnt"
 * @return
 */
func (s *MenuService) GetMenuList(pageDto page_dto.SystemMenuListReqDto, platformType string) (resultDto dto.RespDataListDto, err error) {
	pageDto.Title = helper.DeleteSpace(pageDto.Title)

	var list interface{}
	count := int64(0)

	switch platformType {
	case "admin_plat":
		list, count, err = base_ar.GetPageMenuList(pageDto, &models.PlatMenu{}, &models.PlatMenuMap{})
	case "admin_mchnt":
		list, count, err = base_ar.GetPageMenuList(pageDto, &models.MchntMenu{}, &models.MchntMenuMap{})
	}

	resultDto.List = list
	resultDto.Total = count
	resultDto.PageSize = pageDto.PageSize
	resultDto.CurrentPage = pageDto.CurrentPage
	return resultDto, err
}

// Save
func (s *MenuService) Save(baseParamDto dto.BaseParamDto, menuDto menu_dto.MenuDto, platformType string) (data map[string]string, err error) {
	if baseParamDto.UnitId == "" {
		err = errors.New("Save():参数错误")
		return
	}
	err = global.GetWriteDb().Transaction(func(tx *gorm.DB) error {
		err = s.doSave(tx, baseParamDto, &menuDto, platformType, menuDto.UnitId)
		if err != nil {
			return err
		}
		err = s.doAsyncMenuMap(tx, baseParamDto, menuDto, platformType)
		return err
	})
	// data = make(map[string]string)
	// data["id"] = menuDto.Id
	return
}
func (s *MenuService) doSave(tx *gorm.DB, baseParamDto dto.BaseParamDto, menuDto *menu_dto.MenuDto, platformType string, unitId string) (err error) {
	menuDto.UnitId = unitId
	err = s.checkMenuDto(menuDto, platformType)
	if err != nil {
		return err
	}

	switch platformType {
	case "admin_plat":
		err = base_ar.SaveMenu(tx, menuDto, &models.PlatMenu{})
	case "admin_mchnt":
		err = base_ar.SaveMenu(tx, menuDto, &models.MchntMenu{})
	}

	return
}
func (s *MenuService) doAsyncMenuMap(tx *gorm.DB, baseParamDto dto.BaseParamDto, menuDto menu_dto.MenuDto, platformType string) (err error) {
	userUnitIds := []string{menuDto.UnitId}
	if menuDto.AsyncToAll == "1" {
		unitIdsMap := make(map[string]bool, 0)
		unitIdsMap, err = s.getAllUnitIds(baseParamDto, menuDto, platformType)
		userUnitIds = helper.GetMapKeys(unitIdsMap)
	}
	if err != nil {
		return
	}
	if len(userUnitIds) == 0 {
		err = errors.New("doAsyncMenuMap():没有找到组织单位")
		return
	}

	switch platformType {
	case "admin_plat":
		err = base_ar.RefreshUnitMenuMap(tx, userUnitIds, menuDto.Id, &models.Plat{}, &models.PlatMenu{}, &models.PlatMenuMap{})
	case "admin_mchnt":
		err = base_ar.RefreshUnitMenuMap(tx, userUnitIds, menuDto.Id, &models.Mchnt{}, &models.MchntMenu{}, &models.MchntMenuMap{})
	}

	return
}

// Del
func (s *MenuService) Del(baseParamDto dto.BaseParamDto, menuDto menu_dto.MenuDto, platformType string) (err error) {
	if menuDto.Id == "" {
		return errors.New("Del():参数错误")
	}
	switch platformType {
	case "admin_plat":
		err = base_ar.DelMenu(menuDto.Id, &models.PlatMenu{})
	case "admin_mchnt":
		err = base_ar.DelMenu(menuDto.Id, &models.MchntMenu{})
	}
	return
}

func (s *MenuService) checkMenuDto(menuDto *menu_dto.MenuDto, platformType string) (err error) {
	menuDto.Title = helper.DeleteSpace(menuDto.Title)
	menuDto.Name = helper.DeleteSpace(menuDto.Name)
	menuDto.Path = helper.DeleteSpace(menuDto.Path)
	menuDto.Redirect = helper.DeleteSpace(menuDto.Redirect)
	menuDto.ActivePath = helper.DeleteSpace(menuDto.ActivePath)
	menuDto.Icon = helper.DeleteSpace(menuDto.Icon)
	menuDto.ExtraIcon = helper.DeleteSpace(menuDto.ExtraIcon)
	menuDto.Auths = helper.DeleteSpace(menuDto.Auths)
	menuDto.Component = helper.DeleteSpace(menuDto.Component)

	if menuDto.Title == "" {
		return errors.New("菜单名称不能为空")
	}
	if menuDto.Path == "" {
		return errors.New("菜单路径不能为空")
	}
	if menuDto.UnitId == "" {
		return errors.New("组织单位ID不能为空")
	}

	err = s.checkMenuRepeat(menuDto.Id, menuDto.UnitId, menuDto.Title, platformType)

	return
}

func (s *MenuService) checkMenuRepeat(id, unitId, title string, platformType string) (err error) {
	menuList := make([]base_model.UnitMenu, 0)
	switch platformType {
	case "admin_plat":
		menuList, err = base_ar.GetMenuListByTitle(title, &models.PlatMenu{})
	case "admin_mchnt":
		menuList, err = base_ar.GetMenuListByTitle(title, &models.MchntMenu{})
	}
	addExists := (len(menuList) > 0 && id == "")
	editExists := (len(menuList) > 0 && id != "" && menuList[0].Id != id)

	if len(menuList) > 2 || editExists || addExists {
		return errors.New("菜单名称重复")
	}

	return
}

// 如果要同步菜单前，获取其他组织单位ids
func (s *MenuService) getAllUnitIds(baseParamDto dto.BaseParamDto, menuDto menu_dto.MenuDto, platformType string) (data map[string]bool, err error) {
	data = make(map[string]bool)

	fields := "id, pid, name"
	tmpUnitList := make([]base_model.Unit, 0)
	switch platformType {
	case "admin_plat":
		tmpUnitList, err = base_ar.GetAllUnit(&models.Plat{}, fields)
	case "admin_mchnt":
		tmpUnitList, err = base_ar.GetAllUnit(&models.Mchnt{}, fields)
	}
	if err != nil {
		return data, err
	}

	for _, item := range tmpUnitList {
		data[item.Id] = true
	}
	return
}

// 系统管理-商户菜单管理-商户单位树形
func (s *MenuService) MchntUnitList(baseParamDto dto.BaseParamDto) (data interface{}, err error) {
	tmpData, err := base_ar.GetAllUnit(&models.Mchnt{}, "*")
	for _, item := range tmpData {
		item.LogoLink, _ = helper.LocalFileSign(baseParamDto.Host, item.Logo)
	}
	data = struct {
		List any `json:"list"`
	}{
		List: tmpData,
	}
	return data, err
}
