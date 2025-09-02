package global

import (
	"math/rand"

	"gorm.io/gorm"
)

// 注册中间件，设置日志全局变量
var WriteDb *gorm.DB
var ReadDb []*gorm.DB

func GetWriteDb() *gorm.DB {
	return WriteDb
}

func SetWriteDb(db *gorm.DB) {
	WriteDb = db
}

func GetReadDb() *gorm.DB {
	readDbSize := len(ReadDb)
	if readDbSize > 0 {
		index := rand.Intn(readDbSize)
		return ReadDb[index]
	}
	return WriteDb
}

func SetReadDb(db *gorm.DB) {
	ReadDb = append(ReadDb, db)
}
