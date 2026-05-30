package monitor

import (
	"WenBeego/apps/admin_plat/models_ar"
	dto "WenBeego/apps/common/dto_vo"
	"WenBeego/apps/common/dto_vo/cron_dto"
	"WenBeego/apps/common/dto_vo/page_dto"
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

func (s *CronService) AddTask(baseParamDto dto.BaseParamDto, data cron_dto.UnitCronDto) error {
	return s.saveAndAddTask(baseParamDto, data)
}
func (s *CronService) UpdateTask(baseParamDto dto.BaseParamDto, data cron_dto.UnitCronDto) error {
	return s.saveAndAddTask(baseParamDto, data)
}

func (s *CronService) StartTask(baseParamDto dto.BaseParamDto, data cron_dto.UnitCronDto) error {
	err := s.UpdateTaskStatus(baseParamDto, data, 1)
	return err
}

func (s *CronService) StopTask(baseParamDto dto.BaseParamDto, data cron_dto.UnitCronDto) error {
	crontab.GetCronManager().Stop()
	err := s.UpdateTaskStatus(baseParamDto, data, 0)
	return err
}
func (s *CronService) DelTask(baseParamDto dto.BaseParamDto, data cron_dto.UnitCronDto) error {
	err := global.GetWriteDb().Transaction(func(tx *gorm.DB) error {
		cronDetail, err := s.platCronAr.GetCronById(baseParamDto.UnitId, data.Id)
		if err != nil || cronDetail.Id == "" {
			err = helper.Ternary(err != nil, err, fmt.Errorf("任务不存在"))
			return err
		}

		err = s.platCronAr.Delete(tx, baseParamDto.UnitId, data.Id)
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

func (s *CronService) saveAndAddTask(baseParamDto dto.BaseParamDto, data cron_dto.UnitCronDto) error {
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
		err := s.checkUnitCronDto(baseParamDto, &data)
		if err != nil {
			return err
		}

		err = global.GetWriteDb().Transaction(func(tx *gorm.DB) error {
			if isCreate {
				err = s.platCronAr.Insert(tx, data)
			} else {
				err = s.platCronAr.Update(tx, data)
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

func (s *CronService) UpdateTaskStatus(baseParamDto dto.BaseParamDto, data cron_dto.UnitCronDto, status int) error {
	cronDetail, err := s.platCronAr.GetCronById(baseParamDto.UnitId, data.Id)
	if err != nil {
		return err
	}
	cronDetail.Status = status
	cronDetail.UpdatedBy = &baseParamDto.UnitUserId
	tmpData := cron_dto.UnitCronDto{
		UnitCron: cronDetail,
	}

	err = global.GetWriteDb().Transaction(func(tx *gorm.DB) error {
		err = s.platCronAr.Update(tx, tmpData)
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

func (s *CronService) GetAvailableCronList() (interface{}, error) {
	data := struct {
		List interface{} `json:"list"`
	}{
		List: crontab_task.GetCronTasks(),
	}
	return data, nil
}

func (s *CronService) GetCronList(reqDto *page_dto.MonitorCronListReqDto) (*dto.RespDataListDto, error) {
	data, count, err := s.platCronAr.GetCronList(*reqDto)
	res := &dto.RespDataListDto{}
	res.List = data
	res.Total = count
	res.PageSize = reqDto.PageSize
	res.CurrentPage = reqDto.CurrentPage

	return res, err
}

func (s *CronService) checkUnitCronDto(baseParamDto dto.BaseParamDto, data *cron_dto.UnitCronDto) error {
	var err error

	isCreate := data.Id == ""
	data.UpdatedBy = &baseParamDto.UnitUserId
	if isCreate {
		data.Id, _ = helper.GetUuid()
		data.CreatedBy = baseParamDto.UnitUserId

		err = s.platCronAr.CheckUnitCronDtoEmpty(data)
		if err != nil {
			return err
		}

		exists, err1 := s.platCronAr.GetCronByNameEn(baseParamDto.UnitId, data.NameEn, "")
		if err1 != nil {
			return err1
		}
		if len(exists) > 0 {
			return fmt.Errorf("任务名称已存在")
		}
	} else {
		err = s.platCronAr.CheckUnitCronDtoEmpty(data)
		if err != nil {
			return err
		}

		exists, err1 := s.platCronAr.GetCronByNameEn(baseParamDto.UnitId, data.NameEn, data.Id)
		if err1 != nil {
			return err1
		}
		if len(exists) > 2 || (len(exists) == 0 && exists[0].Id != data.Id) {
			return fmt.Errorf("任务名称已存在")
		}
	}
	return err
}
