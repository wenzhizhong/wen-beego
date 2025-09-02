package helper

import (
	"errors"

	"gorm.io/gorm"
)

// 数据库记录不存在
func DbNotFound(err error) bool {
	return errors.Is(err, gorm.ErrRecordNotFound)
}
