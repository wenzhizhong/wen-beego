package helper

import (
	"errors"
	"math/rand"
	"regexp"
	"strconv"
	"strings"
)

var minPwdLen = 8
var maxPwdLen = 20

func GetSpecialCharacters() string {
	return "!@#$%^&*()_+-=[]{};:,./<>?"
}

// 获取随机密码
func GetRandomPassword(length int) (string, error) {
	if length < minPwdLen {
		return "", errors.New("密码长度不能小于" + strconv.Itoa(minPwdLen))
	}
	if length > maxPwdLen {
		return "", errors.New("密码长度不能大于" + strconv.Itoa(maxPwdLen))
	}
	letterStr := "abcdefghijklmnopqrstuvwxyz"
	letterArr := strings.Split(letterStr, "")
	specialArr := strings.Split(GetSpecialCharacters(), "")
	letterArrLen := len(letterArr)
	specialArrLen := len(specialArr)

	totalPartLen := 4
	itemPartLen := length / totalPartLen
	numPartLen := length - (totalPartLen-1)*itemPartLen

	str := make([]string, length)
	j := 0
	for i := 0; i < itemPartLen; i++ {
		str[j] = letterArr[rand.Intn(letterArrLen)]
		j++
		str[j] = specialArr[rand.Intn(specialArrLen)]
		j++
		str[j] = strings.ToUpper(letterArr[rand.Intn(specialArrLen)])
		j++
	}
	for i := 0; i < numPartLen; i++ {
		str[j] = strconv.Itoa(rand.Intn(10))
		j++
	}

	rand.Shuffle(len(str), func(i, j int) {
		str[i], str[j] = str[j], str[i] // 交换元素
	})
	return strings.Join(str, ""), nil
}
func ValidatePassword(password string) error {
	// 校验长度
	length := len(password)
	if length < minPwdLen {
		return errors.New("密码长度不能小于" + strconv.Itoa(minPwdLen))
	}
	if length > maxPwdLen {
		return errors.New("密码长度不能大于" + strconv.Itoa(maxPwdLen))
	}
	// 检查是否包含非法字符（除了字母、数字、指定特殊字符之外的字符）
	specialChars := regexp.QuoteMeta(GetSpecialCharacters())
	errMsg := "密码只能包含字母、数字和指定特殊字符“" + specialChars + "”"
	if match, _ := regexp.MatchString("[^a-zA-Z0-9"+specialChars+"]", password); match {
		// if match, _ := regexp.MatchString(`[^a-zA-Z0-9!@#$%^&*()_+-=[]{};:,./<>?]`, password); match {
		return errors.New(errMsg + ".")
	}

	// 校验是否包含数字
	if match, _ := regexp.MatchString(`[\d]`, password); !match {
		return errors.New(errMsg + "；")
	}

	// 校验是否包含大小写字母
	hasLower := regexp.MustCompile(`[a-z]`).MatchString(password)
	hasUpper := regexp.MustCompile(`[A-Z]`).MatchString(password)

	if !hasLower || !hasUpper {
		return errors.New(errMsg + "。")
	}

	// 校验是否包含特殊字符
	if match, _ := regexp.MatchString(`[!@#$%^&*()_+\-=\[\]{};:,./<>?]`, password); !match {
		return errors.New(errMsg + "！")
	}
	// 密码强度校验
	err := validatePasswordSafe(password)
	if err != nil {
		return err
	}

	return nil
}

// 密码强度校验，顺序字符(相邻字符ascii相差1)超过5个，则返回弱密码错误
func validatePasswordSafe(password string) error {
	passwordByte := []byte(password)
	passwordLen := len(passwordByte)

	nearly := 0
	nearlyMax := 5
	beginAscii := int('0')
	endAscii := int('z')
	for i := 0; i < passwordLen-1; i++ {
		nearly = 0
		curAscii := int(passwordByte[i])
		if curAscii < beginAscii || curAscii > endAscii {
			continue
		}

		nextAscii := 0
		for j := i + 1; j < passwordLen; j++ {
			curAscii++
			nextAscii = int(passwordByte[j])
			if curAscii == nextAscii {
				if nearly == 0 {
					nearly++
				}
				nearly++
			} else {
				i = j - 1
				break
			}
		}
		if nearly >= nearlyMax && nearlyMax > 0 {
			break
		}
	}

	percent := float32(nearly) / float32(passwordLen)
	if percent >= 0.5 {
		return errors.New("密码过于简单")
	}
	return nil

}
