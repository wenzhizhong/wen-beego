package helper

import (
	"WenBeego/apps/common/dto"
	"WenBeego/apps/common/global"
	"errors"
	"fmt"
	"regexp"
	"time"

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

// 获取验证码
func GetCaptcha(cpatchaType string) (id string, b64s string, answer string, err error) {
	driver, err := getCaptchaDriver(cpatchaType)
	if err != nil {
		return
	}

	catpcha := base64Captcha.NewCaptcha(driver, base64Captcha.DefaultMemStore)
	id, b64s, answer, err = catpcha.Generate()
	return
}

// 获取验证码驱动
func getCaptchaDriver(cpatchaType string) (driver base64Captcha.Driver, err error) {
	switch cpatchaType {
	case dto.AuthCodeTypeDigit:
		driver = &base64Captcha.DriverDigit{
			Height:   80,
			Width:    240,
			Length:   4,
			MaxSkew:  0.7,
			DotCount: 80,
		}
	case dto.AuthCodeTypeMath:
		driver = &base64Captcha.DriverMath{
			Height:          80,
			Width:           240,
			NoiseCount:      0,
			ShowLineOptions: 2 | 4,
			// Length:   6,
			// Source:   "1234567890",
			Fonts: []string{"wqy-microhei.ttc"},
		}
	case dto.AuthCodeTypeString:
		driver = &base64Captcha.DriverString{
			Height:          80,
			Width:           240,
			NoiseCount:      0,
			ShowLineOptions: 2 | 4,
			Length:          6,
			Source:          "1234567890qwertyuiopasdfghjklzxcvbnmQWERTYUIOPASDFGHJKLZXCVBNM",
			Fonts:           []string{"wqy-microhei.ttc"},
		}
	case dto.AuthCodeTypeChinese:
		driver = &base64Captcha.DriverString{
			Height:          80,
			Width:           240,
			NoiseCount:      0,
			ShowLineOptions: 2 | 4,
			Length:          6,
			Source:          "1234567890qwertyuiopasdfghjklzxcvbnmQWERTYUIOPASDFGHJKLZXCVBNM",
		}
	case dto.AuthCodeTypeEmail:
		driver = &base64Captcha.DriverString{
			Height:          80,
			Width:           240,
			NoiseCount:      0,
			ShowLineOptions: 2 | 4,
			Length:          6,
			Source:          "1234567890qwertyuiopasdfghjklzxcvbnmQWERTYUIOPASDFGHJKLZXCVBNM",
		}
	case dto.AuthCodeTypeSms:
		driver = &base64Captcha.DriverString{
			Height:          80,
			Width:           240,
			NoiseCount:      0,
			ShowLineOptions: 2 | 4,
			Length:          6,
			Source:          "1234567890qwertyuiopasdfghjklzxcvbnmQWERTYUIOPASDFGHJKLZXCVBNM",
		}
	default:
		err = errors.New("验证码类型错误")
	}
	return
}

// 校验验证码
func VerifyCaptcha(cpatchaType string, id string, answer string) bool {
	return base64Captcha.DefaultMemStore.Verify(id, answer, true)
}

// 获取时间戳
func GetTimestamp() int64 {
	timezone, _ := global.GetConfigDiy("timezone")
	loc, _ := time.LoadLocation(timezone.(string))
	return time.Now().In(loc).Unix()
}
