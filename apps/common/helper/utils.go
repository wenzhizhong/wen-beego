package helper

import (
	"WenBeego/apps/common/global"
	"WenBeego/apps/common/middleware/captcha_store"
	"WenBeego/apps/common/models"
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"html/template"
	"regexp"
	"strings"
	"time"

	"WenBeego/apps/common/dto"

	googleUuid "github.com/google/uuid"
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
	return regexp.MustCompile(`^1[3-9]\d{9}$`).MatchString(phone)
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

	store := captcha_store.Base64CaptchaRedisStore{}
	store.Expiration = 300 * time.Second
	catpcha := base64Captcha.NewCaptcha(driver, &store)
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
	store := captcha_store.Base64CaptchaRedisStore{}
	return store.Verify(id, answer, true)
}

// 获取时间戳
func GetTimestamp() int64 {
	timezone, _ := global.GetConfigDiy("timezone")
	loc, _ := time.LoadLocation(timezone.(string))
	return time.Now().In(loc).Unix()
}

// 获取时间
func GetTime() time.Time {
	timezone, _ := global.GetConfigDiy("timezone")
	loc, _ := time.LoadLocation(timezone.(string))
	return time.Now().In(loc)
}

// 解析字符串模板
func ParseStringTpl(tpl string, data any) (str string, err error) {
	tmpl, err := template.New("").Parse(tpl)
	if err != nil {
		return str, err
	}
	var result strings.Builder
	err = tmpl.Execute(&result, data)
	if err != nil {
		return str, err
	}

	return result.String(), nil
}

// 判断是否是管理员

func IsAdmin(moduleName string, unitId string, userId string) bool {
	if moduleName == "admin_plat" {
		return getAdminData(unitId, userId, &models.PlatRoleClassify{}, &models.PlatRole{}, &models.PlatUserRole{})
	} else {
		return getAdminData(unitId, userId, &models.MchntRoleClassify{}, &models.MchntRole{}, &models.MchntUserRole{})
	}
}

// 获取管理员用户
func getAdminData[RoleClassifyModel interface{ TableName() string }, RoleModel interface{ TableName() string }, UserRoleModel interface{ TableName() string }](unitId string, userId string, roleClassify RoleClassifyModel, role RoleModel, userRoleModel UserRoleModel) bool {
	tableClassify := roleClassify.TableName()
	tableRole := role.TableName()
	tableUserRole := userRoleModel.TableName()

	err := global.GetReadDb().
		Model(roleClassify).
		Select(tableClassify+".*").
		Joins("inner join "+tableRole+" on "+tableRole+".id = "+tableClassify+".role_id").
		Joins("inner join "+tableUserRole+" on "+tableUserRole+".role_id = "+tableClassify+".role_id").
		Where(tableClassify+".name = ?", "admin").
		Where(tableUserRole+".user_id = ?", userId).
		Where(tableUserRole+".unit_id = ?", unitId).
		Where(tableUserRole+".deleted = ?", 0).
		Where(tableRole+".unit_id = ?", unitId).
		Where(tableRole+".status = ?", 1).
		Where(tableRole+".deleted = ?", 0).
		Take(roleClassify).
		Error
	return err == nil
}

// 获取uuid
func GetUuid() (string, error) {
	googleUuid.EnableRandPool()
	uuid, err := googleUuid.NewV7()
	if err != nil {
		return "", err
	}
	return uuid.String(), nil
}

// 获取md5值
func Md5(str string) string {
	if str == "" {
		return ""
	}

	handle := md5.New()
	handle.Write([]byte(str))
	return hex.EncodeToString(handle.Sum(nil))
}

// 获取refresh token 缓存key
func getRefreshTokenKey(brancaToken string, refreshToken string) string {
	return "refreshToken:" + Md5(refreshToken+":"+brancaToken)
}

// 获取refresh token
func GetRefreshToken(moduleName string, brancaToken string, userId string) (string, error) {
	refresh_token, err := GetUuid()
	if err != nil {
		return "", err
	}
	redisKey := getRefreshTokenKey(brancaToken, refresh_token)

	exp, err := global.GetConfigDiy("branca." + moduleName + ".exp")
	if err != nil {
		return "", err
	}
	refresh_exp := exp.(int) * 2
	if refresh_exp <= 86400 {
		refresh_exp = 2 * 24 * 60 * 60
	}
	err = RedisPut(redisKey, userId, time.Duration(refresh_exp))
	if err != nil {
		return "", err
	}
	return refresh_token, nil
}

// 验证refresh token
func VerifyRefreshToken(brancaToken string, refreshToken string) (result bool, userId string, err error) {
	redisKey := getRefreshTokenKey(brancaToken, refreshToken)
	exits, err := RedisGet(redisKey)
	if err == nil && exits != nil {
		if bytes, ok := exits.([]byte); ok {
			userId = string(bytes)
			DelRefreshToken(brancaToken, refreshToken)
			return true, userId, nil
		}
	}
	return
}

// 删除refresh token
func DelRefreshToken(brancaToken string, refreshToken string) {
	if brancaToken != "" && refreshToken != "" {
		redisKey := getRefreshTokenKey(brancaToken, refreshToken)
		RedisDel(redisKey)
	}
}
