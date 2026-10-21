package base_ar

import (
	"WenBeego/apps/common/dto_vo/menu_dto"
	"WenBeego/apps/common/dto_vo/page_dto"
	"WenBeego/apps/common/global"
	"WenBeego/apps/common/helper"
	"WenBeego/apps/common/models/base_model"
	"WenBeego/apps/common/models/itf"
	"errors"

	"gorm.io/gorm"
)

func GetMenuListByTitle[UnitMenuModel itf.MenuItf](title string, unitMenuModel UnitMenuModel) (data []base_model.UnitMenu, err error) {
	data = make([]base_model.UnitMenu, 0)
	err = global.GetWriteDb().
		Model(unitMenuModel).
		Select("id, title").
		Where("title = ?", title).
		Where("deleted = ?", 0).
		Find(&data).Error
	return
}

// 平台-系统管理-获取菜单列表
func GetPageMenuList[UnitMenuModel itf.MenuItf, UnitMenuMapModel itf.MenuMapItf](pageDto page_dto.SystemMenuListReqDto, unitMenuModel UnitMenuModel, unitMenuMapModel UnitMenuMapModel) (list []base_model.UnitMenu, count int64, err error) {
	list = make([]base_model.UnitMenu, 0)

	tableMenuName := unitMenuModel.TableName()
	tableMenuMapName := unitMenuMapModel.TableName()

	query := global.GetReadDb().Model(unitMenuModel).
		Joins("INNER JOIN "+tableMenuMapName+" ON "+tableMenuMapName+".menu_id = "+tableMenuName+".id").
		Where(tableMenuName+".deleted = ?", 0).
		Where(tableMenuMapName+".deleted = ?", 0)

	if pageDto.Title != "" {
		query = query.Where(tableMenuName+".title LIKE ?", "%"+pageDto.Title+"%")
	}
	if len(pageDto.SelectUnitIds) > 0 {
		query = query.Where(tableMenuMapName+".unit_id IN ?", pageDto.SelectUnitIds)
	}

	err = query.Count(&count).Error
	if err != nil || count == 0 {
		return
	}

	err = query.
		Select(tableMenuName + ".*, " + tableMenuMapName + ".unit_id").
		Find(&list).Error
	return

}

// Save
func SaveMenu[UnitMenuModel itf.MenuItf](tx *gorm.DB, menuDto *menu_dto.MenuDto, unitMenuModel UnitMenuModel) (err error) {
	unitMenuModel.SetId(menuDto.Id)
	unitMenuModel.SetParentId(menuDto.ParentId)
	unitMenuModel.SetUnitId(menuDto.UnitId)
	unitMenuModel.SetMenuType(menuDto.MenuType)
	unitMenuModel.SetTitle(menuDto.Title)
	unitMenuModel.SetName(menuDto.Name)
	unitMenuModel.SetPath(menuDto.Path)
	unitMenuModel.SetComponent(menuDto.Component)
	unitMenuModel.SetRank(menuDto.Rank)
	unitMenuModel.SetRedirect(menuDto.Redirect)
	unitMenuModel.SetIcon(menuDto.Icon)
	unitMenuModel.SetExtraIcon(menuDto.ExtraIcon)
	unitMenuModel.SetEnterTransition(menuDto.EnterTransition)
	unitMenuModel.SetLeaveTransition(menuDto.LeaveTransition)
	unitMenuModel.SetActivePath(menuDto.ActivePath)
	unitMenuModel.SetAuths(menuDto.Auths)
	unitMenuModel.SetFrameSrc(menuDto.FrameSrc)
	unitMenuModel.SetFrameLoading(menuDto.FrameLoading)
	unitMenuModel.SetKeepAlive(menuDto.KeepAlive)
	unitMenuModel.SetHiddenTag(menuDto.HiddenTag)
	unitMenuModel.SetFixedTag(menuDto.FixedTag)
	unitMenuModel.SetShowLink(menuDto.ShowLink)
	unitMenuModel.SetShowParent(menuDto.ShowParent)
	unitMenuModel.SetCreatedAt(menuDto.CreatedAt)
	unitMenuModel.SetUpdatedAt(menuDto.UpdatedAt)
	unitMenuModel.SetDeleted(menuDto.Deleted)
	unitMenuModel.SetClone(menuDto.Clone)
	unitMenuModel.SetMenuFrom(menuDto.MenuFrom)

	if menuDto.Id != "" {
		return EditMenu(tx, unitMenuModel)
	} else {
		menuId, err1 := AddMenu(tx, unitMenuModel)
		if err1 != nil {
			err = err1
		}
		menuDto.Id = menuId
		unitMenuModel.SetId(menuId)
		return
	}
}
func AddMenu[UnitMenuModel itf.MenuItf](tx *gorm.DB, unitMenuModel UnitMenuModel) (menuId string, err error) {
	menuId, err1 := helper.GetUuid()
	if err1 != nil {
		return "", err1
	}
	unitMenuModel.SetId(menuId)
	err = tx.Model(unitMenuModel).Select("*").Create(&unitMenuModel).Error
	return
}
func EditMenu[UnitMenuModel itf.MenuItf](tx *gorm.DB, unitMenuModel UnitMenuModel) (err error) {
	err = tx.Model(unitMenuModel).
		Select("*").
		Omit("deleted", "remark", "clone", "created_at", "created_by").
		Where("id = ?", unitMenuModel.GetId()).
		Updates(&unitMenuModel).Error
	return
}

func DelMenu[UnitMenuModel itf.MenuItf](id string, unitMenuModel UnitMenuModel) (err error) {
	if id == "" {
		return errors.New("Del(): id 参数错误")
	}
	err = global.GetWriteDb().Model(unitMenuModel).Where("id = ?", id).Update("deleted", 1).Error
	return
}
