package helper

import (
	"WenBeego/apps/common/global"
	"WenBeego/apps/common/models"
	"WenBeego/apps/common/models/base_model"
	"WenBeego/apps/common/models/itf"
	"crypto/md5"
	"encoding/hex"
	"errors"
	"html/template"
	"regexp"
	"strings"

	googleUuid "github.com/google/uuid"
)

// Ternary 泛型三元运算符
func Ternary[T any](condition bool, trueVal, falseVal T) T {
	if condition {
		return trueVal
	}
	return falseVal
}

// 指针版本（避免值复制）
func TernaryPtr[T any](condition bool, trueVal, falseVal *T) *T {
	if condition {
		return trueVal
	}
	return falseVal
}

// 字符串专用
func StringTernary(condition bool, trueVal, falseVal string) string {
	if condition {
		return trueVal
	}
	return falseVal
}

// 整数专用
func IntTernary(condition bool, trueVal, falseVal int) int {
	if condition {
		return trueVal
	}
	return falseVal
}

// 布尔专用
func BoolTernary(condition bool, trueVal, falseVal bool) bool {
	if condition {
		return trueVal
	}
	return falseVal
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
	return regexp.MustCompile(`^[\x{4e00}-\x{9fa5}]+$`).MatchString(str)
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

// 判断是否是官方平台id
func IsOfficial(moduleName, unitId string) (bool, error) {
	switch moduleName {
	case "admin_plat":
		platModel := models.Plat{}
		err := global.GetReadDb().Where("id = ?", unitId).Take(&platModel).Error
		if err == nil {
			return platModel.IsOfficial, nil
		}
		if !DbNotFound(err) {
			return false, nil
		}
		return false, err
	default:
		return false, nil
	}
}

// 判断是否是管理员
func IsAdmin(moduleName string, unitUserId string) (bool, error) {
	switch moduleName {
	case "admin_plat":
		return getAdminData(unitUserId, &models.PlatRoleClassify{}, &models.PlatRole{}, &models.PlatUserRole{})
	case "admin_mchnt":
		return getAdminData(unitUserId, &models.MchntRoleClassify{}, &models.MchntRole{}, &models.MchntUserRole{})
	default:
		return false, errors.New("IsAdmin():模块名称错误")
	}
}
func CheckUserHasUnit(moduleName, userId string, requiredUnitIds []string) (bool, error) {
	unitIdsLen := len(requiredUnitIds)
	if unitIdsLen == 0 {
		return false, errors.New("CheckUserHasUnit(): requiredUnitIds 不能为空")
	}
	var unitModel itf.UnitItf
	var unitUserModel itf.UnitUserItf
	tableUnitName := ""
	tableUnitUserName := ""
	dataList := make([]base_model.Unit, 0)

	switch moduleName {
	case "admin_plat":
		unitModel = &models.Plat{}
		unitUserModel = &models.PlatUser{}
	case "admin_mchnt":
		unitModel = &models.Mchnt{}
		unitUserModel = &models.MchntUser{}
	default:
		return false, errors.New("CheckUserHasUnit(): 位置模块" + moduleName)
	}
	tableUnitName = unitModel.TableName()
	tableUnitUserName = unitUserModel.TableName()

	err := global.GetReadDb().
		Model(unitModel).
		Select(tableUnitName+".id").
		Joins("inner join "+tableUnitUserName+" on "+tableUnitUserName+".unit_id = "+tableUnitName+".id").
		Where(tableUnitUserName+".user_id = ?", userId).
		Where(tableUnitUserName+".deleted = ?", 0).
		Where(tableUnitName+".deleted = ?", 0).
		Find(&dataList).Error
	if err != nil || len(dataList) == 0 {
		return false, err
	}

	dataListMap := make(map[string]bool)
	for _, v := range dataList {
		dataListMap[v.Id] = true
	}
	notExistsUnit := make([]string, 0)
	for _, v := range requiredUnitIds {
		if _, ok := dataListMap[v]; !ok {
			notExistsUnit = append(notExistsUnit, v)
		}
	}
	if len(notExistsUnit) > 0 {
		return false, errors.New("您没有以下项目权限：" + strings.Join(notExistsUnit, "、"))
	}

	return true, nil
}

// 获取管理员用户
func getAdminData[RoleClassifyModel itf.RoleClassifyItf, RoleModel itf.RoleItf, UserRoleModel itf.UserRoleItf](unitUserId string, roleClassify RoleClassifyModel, role RoleModel, userRoleModel UserRoleModel) (bool, error) {
	tableClassify := roleClassify.TableName()
	tableRole := role.TableName()
	tableUserRole := userRoleModel.TableName()

	err := global.GetReadDb().
		Model(roleClassify).
		Select(tableClassify+".*").
		Joins("inner join "+tableRole+" on "+tableRole+".id = "+tableClassify+".role_id").
		Joins("inner join "+tableUserRole+" on "+tableUserRole+".role_id = "+tableClassify+".role_id").
		Where(tableClassify+".name = ?", "admin").
		Where(tableUserRole+".user_id = ?", unitUserId).
		Where(tableUserRole+".deleted = ?", 0).
		Where(tableRole+".status = ?", 1).
		Where(tableRole+".deleted = ?", 0).
		Take(roleClassify).
		Error
	if err != nil && !DbNotFound(err) {
		return false, err
	}
	return roleClassify.GetId() != "", nil
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
	err = RedisPut(redisKey, userId, refresh_exp)
	if err != nil {
		return "", err
	}
	return refresh_token, nil
}

// 验证refresh token
func VerifyRefreshToken(brancaToken string, refreshToken string) (result bool, userId string, err error) {
	redisKey := getRefreshTokenKey(brancaToken, refreshToken)
	exits, err := RedisGet(redisKey)
	if err == nil && exits != "" {
		userId = exits
		DelRefreshToken(brancaToken, refreshToken)
		return true, userId, nil
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
