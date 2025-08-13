package global

import (
	"database/sql"
	"math/rand"
)

// 注册中间件，设置日志全局变量
var WriteDb *sql.DB
var ReadDb []*sql.DB

func GetWriteDb() *sql.DB {
	return WriteDb
}

func SetWriteDb(db *sql.DB) {
	WriteDb = db
}

func GetReadDb() *sql.DB {
	readDbSize := len(ReadDb)
	if readDbSize > 0 {
		index := rand.Intn(readDbSize)
		return ReadDb[index]
	}
	return WriteDb
}

func SetReadDb(db *sql.DB) {
	ReadDb = append(ReadDb, db)
}
