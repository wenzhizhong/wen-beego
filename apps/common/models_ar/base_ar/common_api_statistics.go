package base_ar

import (
	"WenBeego/apps/common/global"
	"WenBeego/apps/common/helper"
	"WenBeego/apps/common/models/base_model"
	"WenBeego/apps/common/models/itf"

	"gorm.io/gorm/clause"
)

// 获取今日数据
func GetTodayData[ApiStatisticsModel itf.ApiStatisticsItf](unitId string, apiStatisticsModel ApiStatisticsModel) (data []*base_model.UnitApiStatistics, err error) {
	err = global.GetReadDb().
		Model(&apiStatisticsModel).
		Where("unit_id = ?", unitId).
		Where("date = ?", helper.GetDateStamp()).
		Find(&data).Error
	return
}

func SaveApiLogData[ApiStatisticsModel itf.ApiStatisticsItf](apiStatisticsModel ApiStatisticsModel, data []*base_model.UnitApiStatistics) (err error) {
	// sql := `INSERT INTO your_table (id, field1, field2)
	//        VALUES (@id, @val1, @val2)
	//        ON CONFLICT (id)
	//        DO UPDATE SET field1 = EXCLUDED.field1, field2 = EXCLUDED.field2`
	// db.Exec(sql,
	// 	sql.Named("id", yourID),
	// 	sql.Named("val1", value1),
	// 	sql.Named("val2", value2),
	// )

	if len(data) <= 0 {
		return nil
	}
	err = global.GetWriteDb().
		Model(&apiStatisticsModel).
		Clauses(clause.OnConflict{
			Columns: []clause.Column{{Name: "id"}},
			// UpdateAll: true,
			DoUpdates: clause.AssignmentColumns([]string{"pv", "uv"}),
		}).Create(data).Error
	return err
}
