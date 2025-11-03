package helper

import (
	"WenBeego/apps/common/global"
	"errors"
	"os"
	"strconv"
	"strings"
)

func PathIsExist(path string) bool {
	_, err := os.Stat(path)
	return err == nil || os.IsExist(err)
}

func MkdirAll(path string) error {
	if !PathIsExist(path) {
		return os.MkdirAll(path, os.ModePerm)
	} else {
		fileInfo, _ := os.Stat(path)
		if !fileInfo.IsDir() {
			return errors.New(path + " 不是一个目录")
		}
	}
	return nil
}

// 本地文件访问签名
func LocalFileSign(host, filePath string) (string, error) {
	if filePath == "" {
		return "", nil
	}
	if strings.HasPrefix(filePath, "http") {
		return filePath, nil
	}
	host = Ternary(strings.HasPrefix(host, "http"), host, "//"+host)
	if !strings.HasPrefix(filePath, "/uploads/private/") {
		return host + filePath, nil
	}

	appName, err0 := AppName()
	iss, err1 := global.GetConfigDiy("branca.common.iss")
	aud, err2 := global.GetConfigDiy("branca.common.aud")
	exp, err3 := global.GetConfigDiy("branca.common.exp")
	if err0 != nil || err1 != nil || err2 != nil || err3 != nil || iss == nil || aud == nil || exp == nil {
		return "", errors.New("获取配置异常")
	}

	time := GetTimestamp()
	expInt, err := Interface2Int64(exp)
	if err != nil {
		return "", err
	}

	data := BrancaData{
		Sub:     "",
		SubUnit: appName,
		Exp:     expInt,
		Iat:     time,
		Iss:     iss.(string),
		Aud:     aud.(string),
		Role:    "",
		Scope:   filePath,
	}
	signStr, err := BrancaEncode(data, "common")
	if err != nil {
		return "", err
	}
	sepector := Ternary(strings.Contains(filePath, "?"), "&", "?")
	signStr = host + filePath + sepector + "sign=" + signStr + "&exp=" + strconv.FormatInt(time+expInt, 10)
	return signStr, nil
}

// 本地文件访问签名校验
func LocalFileSignCheck(sign string) (bool, error) {
	if sign == "" {
		return true, nil
	}

	_, err := BrancaDecode(sign, "common")
	if err != nil {
		return false, err
	}
	return true, nil
}
