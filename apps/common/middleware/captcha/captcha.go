package captcha

import (
	"WenBeego/apps/common/dto/auth_dto"
	"errors"
	"time"

	"github.com/mojocn/base64Captcha"
)

// 获取验证码
func GetCaptcha(cpatchaType string) (id string, b64s string, answer string, err error) {
	driver, err := getCaptchaDriver(cpatchaType)
	if err != nil {
		return
	}

	store := Base64CaptchaRedisStore{}
	store.Expiration = 300 * time.Second
	catpcha := base64Captcha.NewCaptcha(driver, &store)
	id, b64s, answer, err = catpcha.Generate()
	return
}

// 获取验证码驱动
func getCaptchaDriver(cpatchaType string) (driver base64Captcha.Driver, err error) {
	switch cpatchaType {
	case auth_dto.AuthCodeTypeDigit:
		driver = &base64Captcha.DriverDigit{
			Height:   80,
			Width:    240,
			Length:   4,
			MaxSkew:  0.7,
			DotCount: 80,
		}
	case auth_dto.AuthCodeTypeMath:
		driver = &base64Captcha.DriverMath{
			Height:          80,
			Width:           240,
			NoiseCount:      0,
			ShowLineOptions: 2 | 4,
			// Length:   6,
			// Source:   "1234567890",
			Fonts: []string{"wqy-microhei.ttc"},
		}
	case auth_dto.AuthCodeTypeString:
		driver = &base64Captcha.DriverString{
			Height:          80,
			Width:           240,
			NoiseCount:      0,
			ShowLineOptions: 2 | 4,
			Length:          6,
			Source:          "1234567890qwertyuiopasdfghjklzxcvbnmQWERTYUIOPASDFGHJKLZXCVBNM",
			Fonts:           []string{"wqy-microhei.ttc"},
		}
	case auth_dto.AuthCodeTypeChinese:
		driver = &base64Captcha.DriverString{
			Height:          80,
			Width:           240,
			NoiseCount:      0,
			ShowLineOptions: 2 | 4,
			Length:          6,
			Source:          "1234567890qwertyuiopasdfghjklzxcvbnmQWERTYUIOPASDFGHJKLZXCVBNM",
		}
	case auth_dto.AuthCodeTypeEmail:
		driver = &base64Captcha.DriverString{
			Height:          80,
			Width:           240,
			NoiseCount:      0,
			ShowLineOptions: 2 | 4,
			Length:          6,
			Source:          "1234567890qwertyuiopasdfghjklzxcvbnmQWERTYUIOPASDFGHJKLZXCVBNM",
		}
	case auth_dto.AuthCodeTypeSms:
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
	store := Base64CaptchaRedisStore{}
	return store.Verify(id, answer, true)
}
