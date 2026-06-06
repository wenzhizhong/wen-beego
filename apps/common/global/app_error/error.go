package app_error

import (
	"WenBeego/apps/common/global/constant"
	"errors"

	"gorm.io/gorm"
)

// 系统错误
func NewSysError(err error, code int) error {
	return baseError(err, code)
}

// 数据库错误
func NewDbError(err error, code ...int) error {
	if err == nil {
		return err
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	tmpCode := constant.ERROR_CODE_DB
	if len(code) > 0 {
		tmpCode = code[0]
	}
	return baseError(err, tmpCode)
}

// 助手错误
func NewHelperError(err error, code ...int) error {
	tmpCode := constant.ERROR_CODE_HELPER
	if len(code) > 0 {
		tmpCode = code[0]
	}
	return baseError(err, tmpCode)
}

// 中间件错误
func NewMiddlewareError(err error, code ...int) error {
	tmpCode := constant.ERROR_CODE_MIDDLEWARE
	if len(code) > 0 {
		tmpCode = code[0]
	}
	return baseError(err, tmpCode)
}
