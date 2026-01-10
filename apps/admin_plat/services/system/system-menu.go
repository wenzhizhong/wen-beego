package system

import (
	"WenBeego/apps/common/dto"
	"WenBeego/apps/common/dto/menu_dto"
	"WenBeego/apps/common/dto/page_dto"
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
		list, count, err = s.commonPlatMenuAr.GetMenuList(pageDto)
	case "admin_mchnt":
		list, count, err = s.commonMchntMenuAr.GetMenuList(pageDto)
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
	global.GetWriteDb().Transaction(func(tx *gorm.DB) error {
		err = s.doSave(tx, baseParamDto, menuDto, platformType, menuDto.UnitId)
		if err != nil {
			return err
		}
		err = s.doSaveOfAsync(tx, baseParamDto, menuDto, platformType)
		return err
	})
	// data = make(map[string]string)
	// data["id"] = menuDto.Id
	return
}
func (s *MenuService) doSave(tx *gorm.DB, baseParamDto dto.BaseParamDto, menuDto menu_dto.MenuDto, platformType string, unitId string) (err error) {
	menuDto.UnitId = unitId
	err = s.checkMenuDto(&menuDto, platformType)
	if err != nil {
		return err
	}

	switch platformType {
	case "admin_plat":
		err = s.commonPlatMenuAr.Save(tx, &menuDto)
	case "admin_mchnt":
		err = s.commonMchntMenuAr.Save(tx, &menuDto)
	}

	return
}
func (s *MenuService) doSaveOfAsync(tx *gorm.DB, baseParamDto dto.BaseParamDto, menuDto menu_dto.MenuDto, platformType string) (err error) {
	if menuDto.AsyncToAll != "1" {
		return
	}
	userUnitIds := make([]string, 0)
	userUnitIds, err = s.getAllUnitIds(baseParamDto, menuDto, platformType)
	if err != nil {
		return
	}

	for _, unitId := range userUnitIds {
		err = s.doSave(tx, baseParamDto, menuDto, platformType, unitId)
		if err != nil {
			return err
		}
	}

	return
}

// Del
func (s *MenuService) Del(baseParamDto dto.BaseParamDto, menuDto menu_dto.MenuDto, platformType string) (err error) {
	if menuDto.Id == "" {
		return errors.New("Del():参数错误")
	}
	menuIds := []string{menuDto.Id}
	userUnitIds := make([]string, 0)
	userUnitIds, err = s.getAllUnitIds(baseParamDto, menuDto, platformType)
	if err != nil {
		return
	}

	meunList := make([]base_model.UnitMenu, 0)
	switch platformType {
	case "admin_plat":
		meunList, err = base_ar.GetMenuListByTitle[*models.PlatMenu](userUnitIds, menuDto.Title)
	case "admin_mchnt":
		meunList, err = base_ar.GetMenuListByTitle[*models.MchntMenu](userUnitIds, menuDto.Title)
	}
	if err != nil {
		return
	}
	for _, v := range meunList {
		if v.Id != menuDto.Id {
			menuIds = append(menuIds, v.Id)
		}
	}

	global.GetWriteDb().Transaction(func(tx *gorm.DB) error {
		for _, menuId := range menuIds {
			switch platformType {
			case "admin_plat":
				err = s.commonPlatMenuAr.Del(menuId)
			case "admin_mchnt":
				err = s.commonMchntMenuAr.Del(menuId)
			}
		}
		return nil
	})
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
		menuList, err = s.commonPlatMenuAr.GetMenuListByTitle(unitId, title)
	case "admin_mchnt":
		menuList, err = s.commonMchntMenuAr.GetMenuListByTitle(unitId, title)
	}
	if len(menuList) > 2 || (len(menuList) > 0 && menuList[0].Id != id) {
		return errors.New("菜单名称重复")
	}

	return
}

// 如果要同步菜单前，获取其他组织单位ids
func (s *MenuService) getAllUnitIds(baseParamDto dto.BaseParamDto, menuDto menu_dto.MenuDto, platformType string) (data []string, err error) {
	data = make([]string, 0)
	userHasUnitIds := make(map[string]bool, 0)

	if menuDto.AsyncToAll == "1" {
		fields := "id, pid, name"
		tmpUnitList := make([]base_model.UnitMenu, 0)
		switch platformType {
		case "admin_plat":
			tmpUnitList, err = base_ar.GetAll(&models.Plat{}, fields)
		case "admin_mchnt":
			tmpUnitList, err = base_ar.GetAll(&models.Mchnt{}, fields)
		}
		if err != nil {
			return data, err
		}

		for _, item := range tmpUnitList {
			if menuDto.UnitId == item.Id {
				continue
			}
			userHasUnitIds[item.Id] = true
		}
	}
	data = helper.GetMapKeys(userHasUnitIds)
	if len(data) == 0 {
		err = errors.New("没有找到用户组织单位列表")
	}
	return
}

func (s *MenuService) GetAllMenu(baseParamDto dto.BaseParamDto, platformType string) (data []base_model.UnitMenu, err error) {
	data = make([]base_model.UnitMenu, 0)
	switch platformType {
	case "admin_plat":

	case "admin_mchnt":

	}
	return
}
