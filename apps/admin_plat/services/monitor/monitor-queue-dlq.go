package monitor

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"WenBeego/apps/common/dto"
	"WenBeego/apps/common/dto/page_dto"
	"WenBeego/apps/common/dto/queue_dlq_dto"
	"WenBeego/apps/common/global"
	"WenBeego/apps/common/global/constant"
	"WenBeego/apps/common/helper"
	"WenBeego/apps/common/middleware/mq"
	"WenBeego/apps/common/models"
	"WenBeego/apps/common/models_ar"
	commonQueueDlq "WenBeego/apps/common/services/queue_dlq"
	"WenBeego/apps/common/thirdPkg/rewrite/RichardKnop/machinery/v1/tasks"

	"gorm.io/gorm"
)

type QueueDlqService struct {
	commonQueueDlq commonQueueDlq.QueueDlq
}

func (s *QueueDlqService) GetList(reqDto page_dto.QueueDlqListReqDto) (dto.RespDataListDto, error) {
	return s.commonQueueDlq.GetList(reqDto)
}

func (s *QueueDlqService) Requeue(baseParamDto dto.BaseParamDto, reqDto queue_dlq_dto.RequeueDto) (int, error) {
	var records []models.QueueDlqFailedLog
	var err error

	if len(reqDto.TaskUUIDs) > 0 {
		records, err = s.commonQueueDlq.GetPendingListByUUIDs(baseParamDto.ModuleName, reqDto.TaskUUIDs)
	} else if reqDto.TaskUUID != "" {
		record, err2 := s.commonQueueDlq.GetByTaskUUID(baseParamDto.ModuleName, reqDto.TaskUUID)
		if err2 != nil {
			return 0, err2
		}
		records = append(records, *record)
	} else {
		condition := page_dto.QueueDlqRequeueReqDto{
			TaskName:        reqDto.TaskName,
			CreateTimeBegin: reqDto.CreateTimeBegin,
			CreateTimeEnd:   reqDto.CreateTimeEnd,
		}
		records, err = s.commonQueueDlq.GetPendingListByCondition(baseParamDto.ModuleName, condition)
	}
	if err != nil {
		return 0, err
	}

	ctx := context.Background()
	successCount := 0

	for _, record := range records {
		lockKey := fmt.Sprintf("queue_dlq_requeue:%s", record.TaskUUID)
		ok, lockErr := global.RedisCache.SetNX(ctx, lockKey, "1", 2*60*time.Second).Result()
		if lockErr != nil || !ok {
			continue
		}

		newUUID, _ := helper.GetUuid()
		// 开启事务
		global.GetWriteDb().Transaction(func(tx *gorm.DB) error {
			insertErr := (&models_ar.QueueDlqFailedRetryAR{}).Insert(tx, &models.QueueDlqFailedRetry{
				Id:    record.TaskUUID,
				NewId: newUUID,
			})
			if insertErr != nil {
				global.Log.Error("Requeue Insert retry error:", insertErr)
				return insertErr
			}
			if updateErr := s.commonQueueDlq.UpdateStatus(tx, baseParamDto.ModuleName, record.TaskUUID, page_dto.QUEUE_DLQ_STATUS_REQUEUED); updateErr != nil {
				global.Log.Error("Requeue UpdateStatus error:", updateErr)
				return updateErr

			}

			var args []tasks.Arg
			if err := json.Unmarshal([]byte(record.TaskArgs), &args); err != nil {
				global.RedisCache.Del(ctx, lockKey)
				return err
			}

			taskName := constant.MqNameType(record.TaskName)
			_, mqErr := (&mq.MqClient{}).SendTask(taskName, args, newUUID)
			if mqErr != nil {
				global.Log.Error("Requeue SendTask error:", mqErr)
				return mqErr
			}

			successCount++
			return nil
		})
	}

	return successCount, nil
}
