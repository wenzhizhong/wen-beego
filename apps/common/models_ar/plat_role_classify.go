package models_ar

import (
	"WenBeego/apps/common/global"
	"WenBeego/apps/common/models"
	"errors"
	"fmt"
)

type PlatRoleClassifyAr struct {
	models.PlatRoleClassify
}

func (a *PlatRoleClassifyAr) GetById(id string) (models.PlatRoleClassify, error) {
	data := models.PlatRoleClassify{}
	if id == "" {
		return data, errors.New("id不能为空")
	}
	result := global.GetReadDb().Where("id = ?", id).Take(&data)
	return data, result.Error
}

/**
 * 通过角色id 获取角色分类
 * @param unitId  组织单位id
 * @param roleId 角色id
 * @return
 */
func (a *PlatRoleClassifyAr) GetClassifysByRoleId(unitId string, roleId string) (platRoleClassify []models.PlatRoleClassify, err error) {
	if unitId == "" || roleId == "" {
		str := fmt.Sprintf("GetClassifysByName():参数不能为空, unitId:%s, roleId:%s", unitId, roleId)
		global.Log.Error(str)
		return platRoleClassify, errors.New(str)
	}
	tableClassify := a.TableName()
	tableRole := (&models.PlatRole{}).TableName()

	err = global.GetReadDb().
		Model(&models.PlatRoleClassify{}).
		Select(tableClassify+".*").
		Where(tableClassify+".role_id = ?", roleId).
		Where(tableRole+".unit_id = ?", unitId).
		Where(tableRole+".status = ?", 1).
		Where(tableRole+".deleted = ?", 0).
		Find(&platRoleClassify).
		Error
	return platRoleClassify, err
}

/**
 * 通过角色名称 获取角色分类
 * @param unitId  组织单位id
 * @param classifyName 分类名称
 * @return
 */
func (a *PlatRoleClassifyAr) GetClassifysByName(unitId string, classifyName string) (platRoleClassify []models.PlatRoleClassify, err error) {
	if unitId == "" || classifyName == "" {
		str := fmt.Sprintf("GetClassifysByName():参数不能为空, unitId:%s, classifyName:%s", unitId, classifyName)
		global.Log.Error(str)
		return platRoleClassify, errors.New(str)
	}
	tableClassify := a.TableName()
	tableRole := (&models.PlatRole{}).TableName()

	err = global.GetReadDb().
		Model(&models.PlatRoleClassify{}).
		Select(tableClassify+".*").
		Where(tableClassify+".name = ?", classifyName).
		Where(tableRole+".unit_id = ?", unitId).
		Where(tableRole+".status = ?", 1).
		Where(tableRole+".deleted = ?", 0).
		Find(&platRoleClassify).
		Error
	return platRoleClassify, err
}
