package base_ar

import (
	"WenBeego/apps/common/dto/user_dto"
	"WenBeego/apps/common/global"
	"WenBeego/apps/common/models"
	"errors"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func SaveUserProfile(tx *gorm.DB, userProfileDto user_dto.UserProfileDto) (err error) {

	userProfileModel := models.UserProfile{}
	err = tx.Model(userProfileModel).Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "id"}},
		DoNothing: true,
		// UpdateAll: true,
		// DoUpdates: clause.AssignmentColumns([]string{
		// 	"avatar",
		// 	"card_type",
		// 	"card_num",
		// 	"card_images",
		// 	"gender",
		// 	"birth_date",
		// 	"constellation",
		// 	"occupation",
		// 	"company",
		// 	"emergency_name",
		// 	"emergency_tel",
		// 	"address",
		// 	"email",
		// 	"source",
		// 	"valid_date_begin",
		// 	"valid_date_end",
		// 	"graduated_from",
		// 	"schooling",
		// 	"degree_number",
		// 	"professional",
		// 	//"status",
		// 	"created_at",
		// 	"updated_at",
		// 	"deleted_at",
		// //	"deleted",
		// 	"remark",
		// }),
	}).Create(&userProfileDto).Error
	return
}

func GetUserProfileOfById(userId string) (models.UserProfile, error) {
	userModel := &models.User{}
	userProfileModel := &models.UserProfile{}

	tableUserName := userModel.TableName()
	tableUserProfileName := userProfileModel.TableName()
	var data models.UserProfile
	if userId == "" {
		return data, errors.New("userId 不能为空")
	}

	result := global.GetReadDb().
		Model(userProfileModel).
		Select(tableUserProfileName+".*").
		Joins("inner join "+tableUserName+" on "+tableUserName+".id = "+tableUserProfileName+".id").
		Where(tableUserName+".user_id = ?", userId).
		Where(tableUserName+".deleted = ?", 0).
		Take(&data)
	return data, result.Error
}
