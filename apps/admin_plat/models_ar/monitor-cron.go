package models_ar

import (
	"WenBeego/apps/common/dto_vo/cron_dto"
	"WenBeego/apps/common/dto_vo/cron_vo"
	"WenBeego/apps/common/dto_vo/page_dto"
	"WenBeego/apps/common/global"
	"WenBeego/apps/common/models"
	"WenBeego/apps/common/models/base_model"
	"fmt"

	"gorm.io/gorm"
)

type PlatCronAr struct {
	models.PlatCron
}

func (ar *PlatCronAr) CheckUnitCronDtoEmpty(data *cron_dto.UnitCronDto) (err error) {
	if data.Id == "" {
		return fmt.Errorf("PlatCronAr Insert(): Id 不能为空")
	}
	if data.Name == "" {
		return fmt.Errorf("PlatCronAr Insert(): Name 不能为空")
	}
	if data.NameEn == "" {
		return fmt.Errorf("PlatCronAr Insert(): NameEn 不能为空")
	}
	if data.CronExpr == "" {
		return fmt.Errorf("PlatCronAr Insert(): CronExpr 不能为空")
	}
	if data.UnitId == "" {
		return fmt.Errorf("PlatCronAr Insert(): CronExpr 不能为空")
	}
	return nil
}
func (ar *PlatCronAr) Insert(tx *gorm.DB, data cron_dto.UnitCronDto) (err error) {
	err = ar.CheckUnitCronDtoEmpty(&data)
	if err != nil {
		return err
	}
	return tx.Model(&models.PlatCron{}).Create(&data).Error
}

func (ar *PlatCronAr) Update(tx *gorm.DB, data cron_dto.UnitCronDto) (err error) {
	err = ar.CheckUnitCronDtoEmpty(&data)
	if err != nil {
		return err
	}
	return global.GetWriteDb().
		Model(&models.PlatCron{}).
		Select("*").
		Omit("unit_id", "created_at", "created_by", "deleted").
		Where("id = ?", data.Id).
		Updates(&data).Error

}

func (ar *PlatCronAr) Delete(tx *gorm.DB, unit_id, id string) (err error) {
	if id == "" {
		return fmt.Errorf("PlatCronAr Delete(): Id 不能为空")
	}
	return tx.Model(&models.PlatCron{}).Where("unit_id = ? AND id = ?", unit_id, id).Update("deleted", 1).Error
}

// 获取计划任务列表
func (ar *PlatCronAr) GetCronList(req page_dto.MonitorCronListReqDto) (data []cron_vo.UnitCronListVo, count int64, err error) {
	data = make([]cron_vo.UnitCronListVo, 0)

	platCronModel := &models.PlatCron{}
	platUserModel := &models.PlatUser{}
	tablePlatCronName := platCronModel.TableName()
	tablePlatUserName := platUserModel.TableName()

	query := global.GetReadDb().
		Model(platCronModel).
		Where(tablePlatCronName+".deleted = ?", 0)

	err = query.Count(&count).Error
	if err != nil {
		return
	}
	if count == 0 {
		return
	}

	err = query.Select(tablePlatCronName + ".*, STRING_AGG(distinct creator.name,'') AS created_by_name, STRING_AGG(distinct updater.name,'') AS updated_by_name").
		Joins("LEFT JOIN " + tablePlatUserName + " AS creator ON creator.id=" + tablePlatCronName + ".created_by").
		Joins("LEFT JOIN " + tablePlatUserName + " AS updater ON updater.id=" + tablePlatCronName + ".updated_by").
		Limit(req.PageSize).
		Offset(req.Offset).
		Group(tablePlatCronName + ".id").
		Order(tablePlatCronName + ".created_at desc").
		Find(&data).Error

	return
}

func (ar *PlatCronAr) GetCronById(unit_id, id string) (data base_model.UnitCron, err error) {
	query := global.GetReadDb().
		Select("*, '' AS created_by_name, '' AS updated_by_name").
		Model(&models.PlatCron{}).
		Where("deleted = 0 AND unit_id = ? AND id = ?", unit_id, id)

	err = query.Take(&data).Error
	return
}
func (ar *PlatCronAr) GetCronByNameEn(unit_id, nameEn, id string) (data []models.PlatCron, err error) {
	query := global.GetReadDb().
		Model(&models.PlatCron{}).
		// Where("deleted = 0 AND \"unit_id\" = ? AND \"name_en\" = ?", unit_id, nameEn)
		Where("deleted = 0 AND \"name_en\" = ?", nameEn)
	if id != "" {
		// query = query.Or("deleted = 0 AND unit_id = ? AND id = ?", unit_id, id)
		query = query.Or("deleted = 0 AND id = ?", id)
	}
	err = query.Find(&data).Error
	return
}
