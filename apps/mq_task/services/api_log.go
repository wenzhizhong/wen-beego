package services

import (
	"WenBeego/apps/common/dto"
)

type ApiLog struct {
}

func (s *ApiLog) SaveToDb(data []dto.ApiLogDto) (interface{}, error) {
	if len(data) <= 0 {
		return nil, nil
	}

	return nil, nil
}

// func doSaveToDb[ApiStatisticsModel any, ApiStatisticsAr any](apiStatisticsModel *ApiStatisticsModel, apiStatisticsAr ApiStatisticsAr, data []dto.ApiLogDto) error {
// 	todayData := []ApiStatisticsModel{}
// 	todayData, err := apiStatisticsAr.GetTodayData()
// 	if err != nil {
// 		return err
// 	}
// 	todayDataMap := map[string]ApiStatisticsModel{}
// 	for _, item := range todayData {
// 		todayDataMap[item.Uri] = item
// 	}

// 	// groupData := map[string][]
// 	// for _, item := range data {
// 	// 	groupData[item.UnitId] = append(groupData[item.UnitId], item)
// 	// }
// }
