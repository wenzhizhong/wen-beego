package middleware

import (
	"WenBeego/apps/common/global"
	"WenBeego/apps/common/helper"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/beego/beego/v2/core/logs"
)

type logConfig struct {
	Filename string   `json:"filename"`
	MaxLines int64    `json:"maxLines"`
	Maxsize  int64    `json:"maxsize"`
	MaxDays  int64    `json:"maxDays"`
	Daily    bool     `json:"daily"`
	Rotate   bool     `json:"rotate"`
	Perm     string   `json:"perm"`
	Separate []string `json:"separate"`
}

func InitLog(logType string) error {
	return initFileLog(logType)
}
func initFileLog(logType string) error {
	fmt.Println("init file log...")

	logConfig, err := getFileLogConfig(logType)
	if err != nil {
		return err
	}
	beeLogger := logs.NewLogger(10000)
	beeLogger.SetLogger(logs.AdapterMultiFile, logConfig)

	// 可选：保留控制台输出
	beeLogger.SetLogger(logs.AdapterConsole)
	// 增强功能
	beeLogger.EnableFuncCallDepth(true)
	beeLogger.SetLogFuncCallDepth(3)
	beeLogger.Async(1e3)

	global.Log = beeLogger

	fmt.Println("init file log done！")
	return nil
}

func getFileLogConfig(logType string) (string, error) {
	mapConfig, err := global.GetConfig("log")
	if err != nil {
		return "", err
	}
	config, _ := helper.MapInterface2MapString(mapConfig)

	defLogPath := global.TempDir + "/logs"
	logPath := config["path"]
	if logType != "" {
		if config[logType] == "" {
			return "", errors.New("log config error, can't find log config for " + logType)
		}
		logPath = config[logType]
	}
	if logPath == "" {
		logPath = defLogPath
	}
	config["path"] = logPath

	// 如果文件夹不存在则创建文件夹
	if !pathExists(logPath) {
		err := os.MkdirAll(logPath, os.ModePerm)
		if err != nil {
			return "", err
		}
	}
	var logConfigObj logConfig
	logConfigObj.Filename = logPath + "/app.log"
	if config["maxlines"] != "" {
		logConfigObj.MaxLines, _ = strconv.ParseInt(config["maxlines"], 10, 64)
	}
	if config["maxsize"] != "" {
		logConfigObj.Maxsize, _ = strconv.ParseInt(config["maxsize"], 10, 64)
	}
	if config["maxdays"] != "" {
		logConfigObj.MaxDays, _ = strconv.ParseInt(config["maxdays"], 10, 64)
	}
	if config["daily"] != "" {
		logConfigObj.Daily, _ = strconv.ParseBool(config["daily"])
	}
	if config["rotate"] != "" {
		logConfigObj.Rotate, _ = strconv.ParseBool(config["rotate"])
	}
	if config["perm"] != "" {
		logConfigObj.Perm = config["perm"]
	}
	logConfigObj.Perm = "0644"
	if config["separate"] != "" {
		logConfigObj.Separate = strings.Split(config["separate"], ",")
	}

	jsonStr, err := json.Marshal(logConfigObj)
	if err != nil {
		return "", err
	}
	return string(jsonStr), nil
}
func pathExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}
