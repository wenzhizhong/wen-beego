package monitor

import (
	"WenBeego/apps/admin_plat/models_ar"
	"WenBeego/apps/common/dto"
	"WenBeego/apps/common/dto/cron_dto"
	"WenBeego/apps/common/dto/page_dto"
	"WenBeego/apps/common/global"
	"WenBeego/apps/common/helper"
	"WenBeego/apps/common/middleware/crontab"
	"WenBeego/apps/common/models"
	"WenBeego/routers/crontab_task"
	"fmt"
	"strings"

	"gorm.io/gorm"
)

type CronService struct {
	CronModel  models.PlatCron
	platCronAr models_ar.PlatCronAr
}

func (c *CronService) AddTask(baseParamDto dto.BaseParamDto, data cron_dto.UnitCronDto) error {
	return c.saveAndAddTask(baseParamDto, data)
}
func (c *CronService) UpdateTask(baseParamDto dto.BaseParamDto, data cron_dto.UnitCronDto) error {
	return c.saveAndAddTask(baseParamDto, data)
}

func (c *CronService) StartTask(baseParamDto dto.BaseParamDto, data cron_dto.UnitCronDto) error {
	err := c.UpdateTaskStatus(baseParamDto, data, 1)
	return err
}

func (c *CronService) StopTask(baseParamDto dto.BaseParamDto, data cron_dto.UnitCronDto) error {
	crontab.GetCronManager().Stop()
	err := c.UpdateTaskStatus(baseParamDto, data, 0)
	return err
}
func (c *CronService) DelTask(baseParamDto dto.BaseParamDto, data cron_dto.UnitCronDto) error {
	err := global.GetWriteDb().Transaction(func(tx *gorm.DB) error {
		cronDetail, err := c.platCronAr.GetCronById(baseParamDto.UnitId, data.Id)
		if err != nil || cronDetail.Id == "" {
			err = helper.Ternary(err != nil, err, fmt.Errorf("任务不存在"))
			return err
		}

		err = c.platCronAr.Delete(tx, baseParamDto.UnitId, data.Id)
		if err != nil {
			return err
		}

		list := crontab_task.GetCronTasks()
		for _, item := range list {
			if item.Name != cronDetail.NameEn {
				continue
			}
			crontab.GetCronManager().RemoveTask(item.Name)
			break
		}
		return nil
	})
	return err
}

func (c *CronService) saveAndAddTask(baseParamDto dto.BaseParamDto, data cron_dto.UnitCronDto) error {
	if data.Name == "" {
		return fmt.Errorf("请选择任务")
	}
	data.UnitId = baseParamDto.UnitId
	data.Group = strings.TrimSpace(data.Group)
	data.Remark = strings.TrimSpace(data.Remark)
	data.CronExpr = strings.TrimSpace(data.CronExpr)

	isCreate := data.Id == ""
	list := crontab_task.GetCronTasks()
	for _, item := range list {
		if item.Name != data.NameEn {
			continue
		}
		err := c.checkUnitCronDto(baseParamDto, &data)
		if err != nil {
			return err
		}

		err = global.GetWriteDb().Transaction(func(tx *gorm.DB) error {
			if isCreate {
				err = c.platCronAr.Insert(tx, data)
			} else {
				err = c.platCronAr.Update(tx, data)
			}
			if err != nil {
				return err
			}
			if data.Status == 0 {
				return nil
			}

			crontab.GetCronManager().RemoveTask(item.Name)
			err = crontab.GetCronManager().AddSafeTask(data.CronExpr, item.CallBack, item.Name)
			return err
		})
		return err
	}
	return nil
}

func (c *CronService) UpdateTaskStatus(baseParamDto dto.BaseParamDto, data cron_dto.UnitCronDto, status int) error {
	cronDetail, err := c.platCronAr.GetCronById(baseParamDto.UnitId, data.Id)
	if err != nil {
		return err
	}
	cronDetail.Status = status
	cronDetail.UpdatedBy = &baseParamDto.UnitUserId
	tmpData := cron_dto.UnitCronDto{
		UnitCron: cronDetail,
	}

	err = global.GetWriteDb().Transaction(func(tx *gorm.DB) error {
		err = c.platCronAr.Update(tx, tmpData)
		if err != nil {
			return err
		}

		list := crontab_task.GetCronTasks()
		for _, item := range list {
			if item.Name != data.NameEn {
				continue
			}

			if status == 1 {
				err = crontab.GetCronManager().AddSafeTask(cronDetail.CronExpr, item.CallBack, item.Name)
			} else {
				crontab.GetCronManager().RemoveTask(item.Name)
			}
			break
		}
		return err
	})
	return err
}

func (c *CronService) GetAvaibleCronList() (interface{}, error) {
	data := struct {
		List interface{} `json:"list"`
	}{
		List: crontab_task.GetCronTasks(),
	}
	return data, nil
}

func (c *CronService) GetCronList(reqDto *page_dto.MonitorCronListReqDto) (*dto.RespDataListDto, error) {
	data, count, err := c.platCronAr.GetCronList(*reqDto)
	res := &dto.RespDataListDto{}
	res.List = data
	res.Total = count
	return res, err
}

func (c *CronService) checkUnitCronDto(baseParamDto dto.BaseParamDto, data *cron_dto.UnitCronDto) error {
	var err error

	isCreate := data.Id == ""
	data.UpdatedBy = &baseParamDto.UnitUserId
	if isCreate {
		data.Id, _ = helper.GetUuid()
		data.CreatedBy = baseParamDto.UnitUserId

		err = c.platCronAr.CheckUnitCronDtoEmpty(data)
		if err != nil {
			return err
		}

		exists, err1 := c.platCronAr.GetCronByNameEn(baseParamDto.UnitId, data.NameEn, "")
		if err1 != nil {
			return err1
		}
		if len(exists) > 0 {
			return fmt.Errorf("任务名称已存在")
		}
	} else {
		err = c.platCronAr.CheckUnitCronDtoEmpty(data)
		if err != nil {
			return err
		}

		exists, err1 := c.platCronAr.GetCronByNameEn(baseParamDto.UnitId, data.NameEn, data.Id)
		if err1 != nil {
			return err1
		}
		if len(exists) > 2 || (len(exists) == 0 && exists[0].Id != data.Id) {
			return fmt.Errorf("任务名称已存在")
		}
	}
	return err
}
