package services

import (
	"WenBeego/apps/common/dto/mq_dto"
	"WenBeego/apps/common/global"
	"WenBeego/apps/common/helper"
	"WenBeego/apps/common/models"
	"WenBeego/apps/common/models/base_model"
	"WenBeego/apps/common/models/itf"
	"WenBeego/apps/common/models_ar/base_ar"
	mqTaskDto "WenBeego/apps/mq_task/dto"
)

type ApiLog struct {
}

func (s *ApiLog) SaveToDb(data []mq_dto.ApiLogDto) (redslt interface{}, err error) {
	if len(data) <= 0 {
		return nil, nil
	}

	//对mq数据进行分组
	groupData, err := s.groupData(data)
	if err != nil {
		return nil, err
	}
	// 新增/更新 api数据统计
	if len(groupData.ModuleMap) <= 0 {
		return nil, nil
	}
	for moduleName, unitGroupData := range groupData.ModuleMap {
		if len(unitGroupData.UnitMap) <= 0 {
			continue
		}
		switch moduleName {
		case "admin_plat":
			apiLogData, err := s.getApiLogSaveData(unitGroupData.UnitMap, &models.PlatMenu{})
			if err != nil {
				return nil, err
			}
			err = base_ar.SaveApiLogData(&models.PlatApiStatistics{}, apiLogData)
			if err != nil {
				return nil, err
			}
		case "admin_mchnt":
			apiLogData, err := s.getApiLogSaveData(unitGroupData.UnitMap, &models.MchntMenu{})
			if err != nil {
				return nil, err
			}
			err = base_ar.SaveApiLogData(&models.MchntApiStatistics{}, apiLogData)
			if err != nil {
				return nil, err
			}
		default:
			global.Log.Warn("SaveToDb：模块不存在: %s", moduleName)
			continue
		}
	}
	return nil, nil
}

func (s *ApiLog) groupData(data []mq_dto.ApiLogDto) (groupData mqTaskDto.ApiLogDataDto, err error) {
	groupData.ModuleMap = make(map[string]*mqTaskDto.ApiLogDataUnitDto)
	for _, item := range data {
		if item.UnitId == "" || item.Uri == "" {
			continue
		}
		moduleName := helper.ParseModuleFromRoute(item.Uri)
		if moduleName == "" {
			continue
		}
		unitMap, ok := groupData.ModuleMap[moduleName]
		if !ok {
			unitMap = &mqTaskDto.ApiLogDataUnitDto{
				UnitMap: make(map[string]*mqTaskDto.ApiLogDataUriDto),
			}
		}
		uriMap, ok := unitMap.UnitMap[item.UnitId]
		if !ok {
			uriMap = &mqTaskDto.ApiLogDataUriDto{
				UriMap: make(map[string][]mq_dto.ApiLogDto),
			}
		}

		uriMap.UriMap[item.Uri] = append(uriMap.UriMap[item.Uri], item)
		unitMap.UnitMap[item.UnitId] = uriMap
		groupData.ModuleMap[moduleName] = unitMap
	}
	return
}

// 获取报错api数据
func (s *ApiLog) getApiLogSaveData(unitMap map[string]*mqTaskDto.ApiLogDataUriDto, menuModel itf.MenuItf) (apiLogData []*base_model.UnitApiStatistics, err error) {
	for unitId, itemUriMap := range unitMap {
		dbPermsMap, err := s.getUnitPermissions(unitId, menuModel)
		if err != nil {
			return apiLogData, err
		}
		todayApiLogMap, err := s.getTodayStatistics(unitId, &models.PlatApiStatistics{})
		if err != nil {
			return apiLogData, err
		}

		for uri, apiDtoArr := range itemUriMap.UriMap {
			uuid, err := helper.GetUuid()
			pv := 0
			uv := 0
			permId := ""
			modulename := helper.ParseModuleFromRoute(uri)
			if err != nil {
				return apiLogData, err
			}
			if dbPermsMap[uri].Id != "" {
				permId = dbPermsMap[uri].Id
			}
			if todayApiLogMap[uri] != (base_model.UnitApiStatistics{}) {
				uuid = todayApiLogMap[uri].ID
				pv = todayApiLogMap[uri].PV
				uv = todayApiLogMap[uri].UV
			}
			for _, apiDto := range apiDtoArr {
				pv++
				if apiDto.UserId != "" {
					uv++
				}
			}

			item := &base_model.UnitApiStatistics{
				ID:         uuid,
				UnitId:     unitId,
				PermsID:    permId,
				URI:        uri,
				PV:         pv,
				UV:         uv,
				Date:       helper.GetDateStamp(),
				Modulename: modulename,
			}
			apiLogData = append(apiLogData, item)
		}
	}
	return apiLogData, nil
}

// 获取组织单位权限列表
func (s *ApiLog) getUnitPermissions(unitId string, menuModel itf.MenuItf) (dbPermsMap map[string]base_model.UnitMenu, err error) {
	dbPermsList, err := base_ar.GetUnitPermissions(unitId, menuModel)
	if err != nil {
		return dbPermsMap, err
	}

	dbPermsMap = make(map[string]base_model.UnitMenu)
	for _, dbPerm := range dbPermsList {
		dbPermsMap[dbPerm.Path] = dbPerm
	}
	return dbPermsMap, nil
}

// 获取今日数据
func (s *ApiLog) getTodayStatistics(unitId string, ApiStatisticsModel itf.ApiStatisticsItf) (dataMap map[string]base_model.UnitApiStatistics, err error) {
	dataList, err := base_ar.GetTodayData(unitId, ApiStatisticsModel)
	if err != nil {
		return dataMap, err
	}

	dataMap = make(map[string]base_model.UnitApiStatistics)
	for _, item := range dataList {
		dataMap[item.URI] = *item
	}
	return dataMap, nil
}
