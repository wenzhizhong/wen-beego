package models_ar

import (
	"WenBeego/apps/admin_plat/models"
	"WenBeego/apps/common/dto_vo/page_dto"
	"WenBeego/apps/common/global"
	"fmt"
	"time"

	"gorm.io/gorm"
)

type SqliteDataSourceAr struct{}

func (ar *SqliteDataSourceAr) Insert(tx *gorm.DB, data *models.SqliteDataSource) error {
	return tx.Model(&models.SqliteDataSource{}).Create(data).Error
}

func (ar *SqliteDataSourceAr) Update(tx *gorm.DB, data *models.SqliteDataSource) error {
	return tx.Model(&models.SqliteDataSource{}).
		Select("name", "file_path", "password", "description", "update_time").
		Where("id = ?", data.Id).
		Updates(data).Error
}

func (ar *SqliteDataSourceAr) Delete(tx *gorm.DB, ids []string) error {
	if len(ids) == 0 {
		return fmt.Errorf("SqliteDataSourceAr Delete(): Ids 不能为空")
	}
	return tx.Model(&models.SqliteDataSource{}).
		Where("id IN ?", ids).
		Updates(map[string]interface{}{"deleted": 1, "update_time": time.Now()}).Error
}

func (ar *SqliteDataSourceAr) GetList(req page_dto.SqliteDataSourceListReqDto) (data []models.SqliteDataSource, count int64, err error) {
	data = make([]models.SqliteDataSource, 0)
	tableName := (&models.SqliteDataSource{}).TableName()

	query := global.GetReadDb().
		Model(&models.SqliteDataSource{}).
		Where(tableName+".deleted = ?", 0)

	if req.Keyword != "" {
		query = query.Where(tableName+".name LIKE ? OR "+tableName+".file_path LIKE ?", "%"+req.Keyword+"%", "%"+req.Keyword+"%")
	}

	err = query.Count(&count).Error
	if err != nil || count == 0 {
		return
	}

	err = query.
		Select(tableName + ".*").
		Limit(req.PageSize).
		Offset(req.Offset).
		Order(tableName + ".create_time desc").
		Find(&data).Error
	return
}

func (ar *SqliteDataSourceAr) GetById(id string) (data models.SqliteDataSource, err error) {
	err = global.GetReadDb().
		Model(&models.SqliteDataSource{}).
		Where("deleted = 0 AND id = ?", id).
		Take(&data).Error
	return
}

func (ar *SqliteDataSourceAr) GetAll() (data []models.SqliteDataSource, err error) {
	data = make([]models.SqliteDataSource, 0)
	err = global.GetReadDb().
		Model(&models.SqliteDataSource{}).
		Where("deleted = 0").
		Order("create_time desc").
		Find(&data).Error
	return
}
