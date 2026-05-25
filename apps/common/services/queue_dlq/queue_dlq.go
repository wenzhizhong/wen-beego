package queue_dlq

import (
	dto "WenBeego/apps/common/dto_vo"
	"WenBeego/apps/common/dto_vo/page_dto"
	"WenBeego/apps/common/dto_vo/queue_dlq_vo"
	"WenBeego/apps/common/helper"
	"WenBeego/apps/common/models"
	"WenBeego/apps/common/models_ar"
	"errors"

	"gorm.io/gorm"
)

type QueueDlq struct{}

func (s *QueueDlq) GetList(reqDto page_dto.QueueDlqListReqDto) (dto.RespDataListDto, error) {
	var data []queue_dlq_vo.QueueDlqListVo
	var count int64
	var err error

	switch reqDto.ModuleName {
	case "admin_plat":
		data, count, err = (&models_ar.QueueDlqFailedLogAR{}).GetList(reqDto)
	default:
		err = errors.New("GetList:模块名称错误")
	}
	if err != nil {
		return dto.RespDataListDto{}, err
	}

	return helper.GetRespDataListDto(reqDto.PageSize, reqDto.CurrentPage, count, data)
}

func (s *QueueDlq) GetPendingListByCondition(moduleName string, reqDto page_dto.QueueDlqRequeueReqDto) ([]models.QueueDlqFailedLog, error) {
	if moduleName != "admin_plat" {
		return nil, errors.New("GetPendingListByCondition:模块名称错误")
	}
	return (&models_ar.QueueDlqFailedLogAR{}).GetPendingListByCondition(reqDto)
}

func (s *QueueDlq) GetPendingListByUUIDs(moduleName string, taskUUIDs []string) ([]models.QueueDlqFailedLog, error) {
	if moduleName != "admin_plat" {
		return nil, errors.New("GetPendingListByUUIDs:模块名称错误")
	}
	return (&models_ar.QueueDlqFailedLogAR{}).GetPendingListByUUIDs(taskUUIDs)
}

func (s *QueueDlq) GetByTaskUUID(moduleName string, taskUUID string) (*models.QueueDlqFailedLog, error) {
	if moduleName != "admin_plat" {
		return nil, errors.New("GetByTaskUUID:模块名称错误")
	}
	return (&models_ar.QueueDlqFailedLogAR{}).GetByTaskUUID(taskUUID)
}

func (s *QueueDlq) UpdateStatus(tx *gorm.DB, moduleName string, taskUUID string, status int) error {
	if moduleName != "admin_plat" {
		return errors.New("UpdateStatus:模块名称错误")
	}
	return (&models_ar.QueueDlqFailedLogAR{}).UpdateStatus(tx, taskUUID, status)
}
