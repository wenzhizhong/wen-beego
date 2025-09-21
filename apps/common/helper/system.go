package helper

import (
	"WenBeego/apps/common/global"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/beego/beego/v2/server/web/context"
)

// ASCII 艺术字符
func Get_ASCII_ArtisticCharacters() (string, error) {
	// static/common/text/logo.txt
	file, err := os.Open("static/common/text/logo.txt")
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

// 从路由路径解析模块名
func ParseModuleFromRoute(ctx *context.Context) string {
	path := ctx.Request.URL.Path

	// 移除开头的斜杠
	if strings.HasPrefix(path, "/") {
		path = path[1:]
	}

	// 分割路径部分
	parts := strings.Split(path, "/")

	// 路径格式: /模块/控制器/方法
	if len(parts) >= 1 {
		return parts[0]
	}

	return "index"
}
