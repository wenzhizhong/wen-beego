package helper

import (
	"WenBeego/apps/common/global"
	"encoding/json"
	"errors"
	"strings"

	"github.com/hako/branca"
)

type BrancaData struct {
	Sub         string `json:"sub" `         // 用户ID（user表，Subject用户主体）
	SubUnit     string `json:"subUnit" `     // 用户所在单位组织
	SubUnitUser string `json:"subUnitUser" ` // 用户所在单位组织的id（plat_user表/mchnt_user表，用户所在单位组织）
	Exp         int64  `json:"exp"`          // 过期时间（Unix 时间戳）
	Iat         int64  `json:"iat"`          // 签发时间（可省略，Header 中已包含时间戳）
	Iss         string `json:"iss"`          // 签发者 (Issuer)
	Aud         string `json:"aud"`          // 接收者 (Audience)
	Role        string `json:"role"`         // 自定义字段（用户角色）
	Scope       string `json:"scope"`        // 自定义字段（权限范围）
}

// encode
func BrancaEncode(data BrancaData, moduleName string) (string, error) {
	key, err := getBranceKey(moduleName)
	if err != nil {
		return "", err
	}

	needEncodeString, err := json.Marshal(data)
	if err != nil {
		return "", err
	}

	brancaObj := branca.NewBranca(key)
	return brancaObj.EncodeToString(string(needEncodeString))
}

// decode
func BrancaDecode(needDecodeString string, moduleName string) (BrancaData, error) {
	needDecodeString = strings.TrimPrefix(needDecodeString, "Bearer ")
	key, err := getBranceKey(moduleName)
	if err != nil {
		return BrancaData{}, err
	}

	brancaObj := branca.NewBranca(key)
	dataStr, err := brancaObj.DecodeToString(needDecodeString)
	if err != nil {
		return BrancaData{}, err
	}

	data := BrancaData{}
	err = json.Unmarshal([]byte(dataStr), &data)
	if err != nil {
		return BrancaData{}, err
	}
	return data, nil
}

// Brance key
func getBranceKey(moduleName string) (string, error) {
	key, err := global.GetConfigDiy("branca." + moduleName + ".key")
	if err != nil || key == nil {
		global.Log.Error("请配置branca key")
		return "", errors.New("请配置branca key")
	}
	return key.(string), nil
}
