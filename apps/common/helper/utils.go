package helper

import (
	"errors"
	"fmt"
	"regexp"

	"github.com/mojocn/base64Captcha"
)

// interface{} 转 map[string]interface{}
func Interface2MapInterface(i interface{}) (map[string]interface{}, error) {
	if i == nil {
		return make(map[string]interface{}), nil
	}
	tmpMapConfig, ok := i.(map[string]interface{})
	if !ok {
		fmt.Println("类型转换错误, \ninput=", i)
		return nil, errors.New("类型转换错误")
	}
	return tmpMapConfig, nil
}

// map[string]interface{} 转 map[string]string
func MapInterface2MapString(i map[string]interface{}) (map[string]string, error) {
	returnMap := make(map[string]string)
	for k, v := range i {
		returnMap[k] = fmt.Sprint(v)
		if v == nil {
			returnMap[k] = ""
		}
	}
	return returnMap, nil
}

// 判断手机号
func IsCellPhone(phone string) bool {
	return regexp.MustCompile(`^1[3-9]\d{8}$`).MatchString(phone)
}

// 判断邮箱
func IsEmail(email string) bool {
	return regexp.MustCompile(`^[a-zA-Z0-9_-]+@[a-zA-Z0-9_-]+(\.[a-zA-Z0-9_-]+)+$`).MatchString(email)
}

// 判断中文
func IsChinese(str string) bool {
	return regexp.MustCompile(`^[\u4e00-\u9fa5]+$`).MatchString(str)
}

func GetCaptcha(cpatchaType string) (id string, b64s string, answer string, err error) {
	var driver = base64Captcha.DriverDigit{
		Height:   80,
		Width:    240,
		Length:   4,
		MaxSkew:  0.7,
		DotCount: 80,
	}
	catpcha := base64Captcha.NewCaptcha(&driver, base64Captcha.DefaultMemStore)

	id, b64s, answer, err = catpcha.Generate()
	return
}

func VerifyCaptcha(cpatchaType string, id string, answer string) bool {
	return base64Captcha.DefaultMemStore.Verify(id, answer, true)
}
