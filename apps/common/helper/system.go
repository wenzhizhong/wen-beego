package helper

import (
	"WenBeego/apps/common/global"
	"errors"
	"fmt"
	"io"
	"os"
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

func Output_ASCII_ArtisticCharacters() {
	logo, _ := Get_ASCII_ArtisticCharacters()
	version, _ := AppVersion()

	fmt.Println(logo + fmt.Sprintf("\nAppVersion: %s\n\n", version))
}
func AppVersion() (string, error) {
	mapConfig, err := global.GetConfig("default")
	if err != nil {
		return "", errors.New("获取版本号错误")
	}
	version := mapConfig["version"].(string)
	return version, nil
}
