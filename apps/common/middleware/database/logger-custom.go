package database

import (
	"log"
	"os"
	"time"

	"gorm.io/gorm/logger"
)

// 创建一个新的自定义 logger 实例
func NewMyGormLogger(config Config) logger.Interface {

	// 使用标准输出构造 Writer
	writer := log.New(os.Stdout, "\r\n", log.LstdFlags)
	defaultLogger := New(writer, config)

	return defaultLogger
}

// 或者创建一个基于特定配置的 logger
func NewMyGormLoggerWithConfig(slowThreshold time.Duration, logLevel logger.LogLevel, ignoreRecordNotFoundError, parameterizedQueries, colorful bool) logger.Interface {
	config := Config{
		SlowThreshold:             slowThreshold,
		LogLevel:                  logLevel,
		IgnoreRecordNotFoundError: ignoreRecordNotFoundError,
		ParameterizedQueries:      parameterizedQueries,
		Colorful:                  colorful,
	}

	writer := log.New(os.Stdout, "\r\n", log.LstdFlags)
	defaultLogger := New(writer, config)

	return defaultLogger
}
