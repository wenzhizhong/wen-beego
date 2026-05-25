package models_ar

import (
	"WenBeego/apps/common/dto_vo/generate_code_dto"
	"WenBeego/apps/common/dto_vo/generate_code_vo"
	"WenBeego/apps/common/dto_vo/page_dto"
	"WenBeego/apps/common/global"
	"WenBeego/apps/common/models"
	"WenBeego/apps/common/models/base_model"
	"fmt"
	"time"

	"gorm.io/gorm"
)

type GenerateCodeAr struct {
	models.GenerateCode
}

func (ar *GenerateCodeAr) Insert(tx *gorm.DB, data generate_code_dto.GenerateCodeDto) error {
	now := time.Now()
	data.CreateTime = &now
	err := tx.Model(&models.GenerateCode{}).Create(&data).Error
	return err
}

func (ar *GenerateCodeAr) Update(tx *gorm.DB, data generate_code_dto.GenerateCodeDto) error {
	return tx.Model(&models.GenerateCode{}).
		Select("*").
		Omit("id", "deleted").
		Where("id = ?", data.Id).
		Updates(&data).Error
}

func (ar *GenerateCodeAr) Delete(tx *gorm.DB, ids []string) error {
	if len(ids) == 0 {
		return fmt.Errorf("GenerateCodeAr Delete(): Ids 不能为空")
	}
	return tx.Model(&models.GenerateCode{}).Where("id IN ?", ids).Update("deleted", 1).Error
}

func (ar *GenerateCodeAr) GetList(req page_dto.GenerateCodeListReqDto) (data []generate_code_vo.GenerateCodeListVo, count int64, err error) {
	data = make([]generate_code_vo.GenerateCodeListVo, 0)
	model := &models.GenerateCode{}
	tableName := model.TableName()

	query := global.GetReadDb().
		Model(model).
		Where(tableName+".deleted = ?", 0)

	if req.Keyword != "" {
		query = query.Where(tableName+".table_name LIKE ?", "%"+req.Keyword+"%")
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

func (ar *GenerateCodeAr) GetById(id string) (data base_model.GenerateCode, err error) {
	err = global.GetReadDb().
		Select("*, '' AS create_by_name, '' AS update_by_name").
		Model(&models.GenerateCode{}).
		Where("deleted = 0 AND id = ?", id).
		Take(&data).Error
	return
}

func (ar *GenerateCodeAr) GetAllDbTables() (data []string, err error) {
	data = make([]string, 0)
	err = global.GetReadDb().Raw(
		"SELECT table_name FROM information_schema.tables WHERE table_schema = 'public' AND table_type = 'BASE TABLE' ORDER BY table_name",
	).Scan(&data).Error
	return
}

type DbColumnInfo struct {
	ColumnName    string  `gorm:"column:column_name" json:"column_name"`
	DataType      string  `gorm:"column:data_type" json:"data_type"`
	CharMaxLength *int    `gorm:"column:character_maximum_length" json:"character_maximum_length"`
	IsNullable    string  `gorm:"column:is_nullable" json:"is_nullable"`
	ColumnDefault *string `gorm:"column:column_default" json:"column_default"`
	Comment       string  `gorm:"column:comment" json:"comment"`
}

func (ar *GenerateCodeAr) GetDbTableColumns(tableName string) (data []DbColumnInfo, err error) {
	data = make([]DbColumnInfo, 0)
	err = global.GetReadDb().Raw(`
		SELECT 
			c.column_name,
			c.data_type,
			c.character_maximum_length,
			c.is_nullable,
			c.column_default,
			COALESCE(col_description(pc.oid, c.ordinal_position), '') AS comment
		FROM information_schema.columns c
		JOIN pg_class pc ON pc.relname = c.table_name
		JOIN pg_namespace pn ON pn.oid = pc.relnamespace AND pn.nspname = c.table_schema
		WHERE c.table_schema = 'public' AND c.table_name = ?
		ORDER BY c.ordinal_position
	`, tableName).Scan(&data).Error
	return
}

func (ar *GenerateCodeAr) GetTableComment(tableName string) (string, error) {
	var comment string
	err := global.GetReadDb().Raw(`
		SELECT COALESCE(obj_description(pc.oid), '') 
		FROM pg_class pc 
		JOIN pg_namespace pn ON pn.oid = pc.relnamespace 
		WHERE pc.relname = ? AND pn.nspname = 'public'
	`, tableName).Scan(&comment).Error
	return comment, err
}
