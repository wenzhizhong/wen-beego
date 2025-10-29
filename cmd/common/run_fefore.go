package common

import (
	"WenBeego/apps/common/global"
	"os"
	"path/filepath"
)

func RunBefore() {
	// 全局目录
	initGlobalPath()
}

func initGlobalPath() {
	rootPath, _ := filepath.Abs("../../")
	appDir := filepath.Join(rootPath, "apps")
	configDir := filepath.Join(rootPath, "conf")
	staticDir := filepath.Join(rootPath, "static")
	routersDir := filepath.Join(rootPath, "routers")
	tempDir := filepath.Join(rootPath, "temp")

	global.RootPath = rootPath
	global.AppDir = appDir
	global.ConfigDir = configDir
	global.StaticDir = staticDir
	global.RoutersDir = routersDir
	global.TempDir = tempDir

	fileInfo, err := os.Stat(tempDir)
	if os.IsNotExist(err) || !fileInfo.IsDir() {
		os.MkdirAll(tempDir, os.ModePerm)
	}

}
