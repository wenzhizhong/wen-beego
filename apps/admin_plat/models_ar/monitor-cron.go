package models_ar

import (
	"WenBeego/apps/common/dto/cron_dto"
	"WenBeego/apps/common/dto/page_dto"
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
	return tx.Create(&data).Error
}

func (ar *PlatCronAr) Update(tx *gorm.DB, data cron_dto.UnitCronDto) (err error) {
	err = ar.CheckUnitCronDtoEmpty(&data)
	if err != nil {
		return err
	}
	return global.GetWriteDb().Model(&data).Omit("created_at", "created_by", "deleted").Updates(&data).Error

}

func (ar *PlatCronAr) Delete(id string) (err error) {
	if id == "" {
		return fmt.Errorf("PlatCronAr Delete(): Id 不能为空")
	}
	return global.GetWriteDb().Model(&models.PlatCron{}).Where("id = ?", id).Update("deleted", 1).Error
}

// 获取计划任务列表
func (ar *PlatCronAr) GetCronList(req page_dto.MonitorCronListReqDto) (data []models.PlatCron, count int64, err error) {
	data = make([]models.PlatCron, 0)
	query := global.GetReadDb().
		Where("deleted = ?", 0)

	err = query.Count(&count).Error
	if err != nil {
		return
	}
	if count == 0 {
		return
	}

	err = query.Limit(req.PageSize).Offset(req.Offset).Find(&data).Error

	return
}

func (ar *PlatCronAr) GetCronById(unit_id, id string) (data base_model.UnitCron, err error) {
	query := global.GetReadDb().
		Model(&models.PlatCron{}).
		Where("deleted = 0 AND unit_id = ? AND id = ?", unit_id, id)

	err = query.Take(&data).Error
	return
}
func (ar *PlatCronAr) GetCronByName(unit_id, name, id string) (data []models.PlatCron, err error) {
	query := global.GetReadDb().
		Model(&models.PlatCron{}).
		Where("deleted = 0 AND \"unit_id\" = ? AND \"name\" = ?", unit_id, name)
	if id != "" {
		query = query.Or("deleted = 0 AND unit_id = ? AND id = ?", unit_id, id)
	}
	err = query.Find(&data).Error
	return
}
