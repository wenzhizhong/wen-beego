package base_ar

import (
	"WenBeego/apps/common/dto_vo/user_dto"
	"WenBeego/apps/common/global"
	"WenBeego/apps/common/models"
	"errors"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func SaveUser(tx *gorm.DB, userDto user_dto.UserDto) (err error) {
	userModel := models.User{}
	err = tx.Model(userModel).Clauses(clause.OnConflict{
		Columns: []clause.Column{{Name: "id"}},
		// UpdateAll: true,
		DoUpdates: clause.AssignmentColumns([]string{"name", "phone", "email"}),
	}).Create(&userDto).Error
	return
}

func GetUserByIdOrPhone(id, phone string) (data []user_dto.UserAllDataDto, err error) {
	if id == "" && phone == "" {
		err = errors.New("GetUserByIdOrPhone(): 参数id，phone不能均为空")
		return
	}
	tableUser := (&models.User{}).TableName()
	tableUserProfile := (&models.UserProfile{}).TableName()

	query := global.GetReadDb().
		Model(models.User{}).
		Select("*").
		Joins("INNER JOIN "+tableUserProfile+" ON "+tableUserProfile+".id = \""+tableUser+"\".id").
		Where(tableUserProfile+".deleted = ?", 0)
	if id != "" {
		query = query.Where("\""+tableUser+"\".id = ?", id)
	} else {
		query = query.Where("\""+tableUser+"\".phone = ?", phone)
	}

	err = query.Find(&data).Error
	for i := range data {
		data[i].User.Id = data[i].Id
		data[i].UserProfile.Id = data[i].Id
	}
	return
}

func GetUserById(id string) (data []user_dto.UserAllDataDto, err error) {
	tableUser := (&models.User{}).TableName()
	tableUserProfile := (&models.UserProfile{}).TableName()

	err = global.GetReadDb().
		Model(models.User{}).
		Select("*").
		Joins("INNER JOIN "+tableUserProfile+" ON "+tableUserProfile+".id = \""+tableUser+"\".id").
		Where("\""+tableUser+"\".id = ?", id).
		Where(tableUserProfile+".deleted = ?", 0).
		Find(&data).Error
	return
}
