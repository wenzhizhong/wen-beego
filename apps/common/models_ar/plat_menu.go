package models_ar

import (
	"WenBeego/apps/common/dto/menu_dto"
	"WenBeego/apps/common/dto/page_dto"
	"WenBeego/apps/common/global"
	"WenBeego/apps/common/helper"
	"WenBeego/apps/common/models"
	"WenBeego/apps/common/models/base_model"
	"errors"

	"gorm.io/gorm"
)

type PlatMenuAr struct {
	models.PlatMenu
}

func (a *PlatMenuAr) GetById(id string) (models.PlatMenu, error) {
	data := models.PlatMenu{}
	if id == "" {
		return data, errors.New("id不能为空")
	}
	result := global.GetReadDb().Where("id = ?", id).Take(&data)
	return data, result.Error
}

// 平台-系统管理-获取菜单列表
func (a *PlatMenuAr) GetMenuList(pageDto page_dto.SystemMenuListReqDto) (list []models.PlatMenu, count int64, err error) {
	list = make([]models.PlatMenu, 0)

	query := global.GetReadDb().Model(&models.PlatMenu{}).Where("deleted = ?", 0)
	if pageDto.Title != "" {
		query = query.Where("title LIKE ?", "%"+pageDto.Title+"%")
	}
	if len(pageDto.SelectUnitIds) > 0 {
		query = query.Where("unit_id IN ?", pageDto.SelectUnitIds)
	}

	err = query.Count(&count).Error
	if err != nil || count == 0 {
		return
	}

	err = query.Find(&list).Error
	return

}

// Save
func (a *PlatMenuAr) Save(tx *gorm.DB, menuDto *menu_dto.MenuDto) (err error) {
	if menuDto.Id != "" {
		return a.Edit(tx, menuDto)
	} else {
		return a.Add(tx, menuDto)
	}
}
func (a *PlatMenuAr) Add(tx *gorm.DB, menuDto *menu_dto.MenuDto) (err error) {
	menuDto.Id, err = helper.GetUuid()
	if err != nil {
		return
	}
	err = tx.Model(&models.PlatMenu{}).Select("*").Create(&menuDto).Error
	return
}
func (a *PlatMenuAr) Edit(tx *gorm.DB, menuDto *menu_dto.MenuDto) (err error) {
	err = tx.Model(&models.PlatMenu{}).
		Select("*").
		Omit("deleted", "remark", "clone", "created_at", "created_by").
		Where("id = ?", menuDto.Id).
		Updates(&menuDto).Error
	return
}

func (a *PlatMenuAr) Del(id string) (err error) {
	if id == "" {
		return errors.New("Del(): id 参数错误")
	}
	err = global.GetWriteDb().Model(&models.PlatMenu{}).Where("id = ?", id).Update("deleted", 1).Error
	return
}

func (a *PlatMenuAr) GetMenuListByTitle(unitId, title string) (data []base_model.UnitMenu, err error) {
	data = make([]base_model.UnitMenu, 0)
	err = global.GetWriteDb().
		Model(&models.PlatMenu{}).
		Select("id, title").
		Where("unit_id = ?", unitId).
		Where("title = ?", title).
		Where("deleted = ?", 0).
		First(&data).Error
	return
}
