package middleware

import (
	"WenBeego/apps/common/global"
	"WenBeego/apps/common/helper"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// 创建数据库连接
func InitDb() error {
	return initPgSql()
}

// pgsql
func initPgSql() error {
	// 获取配置
	writeDbKey := "pgsql"
	readDbKeys := "pgsql_read"
	config, err := getPgConfig(writeDbKey)
	if err != nil {
		return err
	}
	readConfigs, err := getReadPgConfigs(readDbKeys)
	if err != nil {
		return err
	}

	// 初始化写入数据库
	fmt.Println("init write db...")
	writeDb, err := doInitPgSql(config)
	if err != nil {
		return err
	}
	global.SetWriteDb(writeDb)

	// 初始化读取数据库
	fmt.Println("init read dbs...")
	for _, readConfig := range readConfigs {
		readDb, err := doInitPgSql(readConfig)
		if err != nil {
			return err
		}
		global.SetReadDb(readDb)
	}

	fmt.Println("init db done！")
	return nil
}
func doInitPgSql(config map[string]string) (*sql.DB, error) {
	database := &sql.DB{}
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=%s", config["host"], config["user"], config["password"], config["dbname"], config["port"], config["sslmode"], config["timezone"])
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return database, err
	}

	database, err = db.DB()
	if err != nil {
		return database, err
	}
	pingErr := database.Ping()
	if pingErr != nil {
		return database, err
	}

	return database, nil
}
func initMySql() {
}

func getPgConfig(keys string) (map[string]string, error) {
	if keys == "" {
		return nil, errors.New("getPgConfig(): param ‘keys’ is empty")
	}
	keysArr := strings.Split(keys, ".")
	config, tmpMapConfig, err := make(map[string]string), make(map[string]interface{}), *new(error)

	// 获取一级配置
	tmpMapConfig, tmpErr := global.GetConfig(keysArr[0])
	if tmpErr != nil {
		return nil, tmpErr
	}
	// 获取深层配置
	if len(keysArr) > 1 {
		for index, key := range keysArr {
			if index == 0 {
				continue
			}
			tmpMapConfig, err = helper.Interface2MapInterface(tmpMapConfig[key])
			if err != nil {
				return nil, err
			}
		}
	}
	config, _ = helper.MapInterface2MapString(tmpMapConfig)

	if !(config["host"] != "" && config["port"] != "" && config["dbname"] != "" && config["user"] != "" && config["password"] != "" && config["sslmode"] != "" && config["timezone"] != "") {
		return nil, errors.New("pgsql config error")
	}
	return config, nil
}
func getReadPgConfigs(firstLevelKey string) ([]map[string]string, error) {
	keys, err := getReadDdKeys(firstLevelKey)
	if err != nil {
		return make([]map[string]string, 0), err
	}

	configs := make([]map[string]string, 0)
	for _, key := range keys {
		tmpKey := firstLevelKey + "." + key
		config, err := getPgConfig(tmpKey)
		if err != nil {
			return make([]map[string]string, 0), err
		}
		configs = append(configs, config)
	}
	return configs, nil
}
func getReadDdKeys(firstLevelKey string) ([]string, error) {
	mapConfigData, tmpErr := global.GetConfig(firstLevelKey)
	if tmpErr != nil {
		return nil, tmpErr
	}

	keys := []string{}
	for key, _ := range mapConfigData {
		keys = append(keys, key)
	}

	return keys, nil
}
