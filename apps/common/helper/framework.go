package helper

import (
	"WenBeego/apps/common/global"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
)

// ASCII 艺术字符
func Get_ASCII_ArtisticCharacters() (string, error) {
	// static/common/text/logo.txt
	file, err := os.Open(global.StaticDir+"/common/text/logo.txt")
	if err != nil {
		return "", errors.New("打开文件错误")
	}
	defer file.Close()

	content, err := io.ReadAll(file)
	if err != nil {
		return "", errors.New("读取文件错误")
	}
	return string(content), nil
}

// 输出艺术字logo到终端
func Output_ASCII_ArtisticCharacters() {
	logo, _ := Get_ASCII_ArtisticCharacters()
	version, _ := AppVersion()

	fmt.Println(logo + fmt.Sprintf("\nAppVersion: %s\n\n", version))
}

// 获取框架版本
func AppVersion() (string, error) {
	tmpVersion, err := global.GetConfigDiy("version")
	if err != nil {
		return "", errors.New("获取版本号错误")
	}
	version := tmpVersion.(string)
	return version, nil
}

// 获取框架runmode
func AppRunmode() (string, error) {
	tmp, err := global.GetConfigDiy("runmode")
	if err != nil {
		return "", errors.New("获取runmode错误")
	}
	data := tmp.(string)
	return data, nil
}

// 检查运行模式
func CheckRunMode(runMode string) bool {
	if runMode, err := AppRunmode(); err == nil && runMode == "dev" {
		return true
	}
	return false
}

// 从路由路径解析模块名
func ParseModuleFromRoute(path string) string {
	path = strings.TrimPrefix(path, "/")
	parts := strings.Split(path, "/")

	// 路径格式: /模块/控制器/方法
	if len(parts) >= 1 {
		return parts[0]
	}

	return "index"
}
