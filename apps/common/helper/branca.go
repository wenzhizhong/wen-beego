package helper

import (
	"WenBeego/apps/common/global"

	"github.com/hako/branca"
)

const defKey = "DdJPX4TcgGBhliONFebwjQu9Y2v1fpxs"

// encode
func BrancaEncode(needEncodeString string) (string, error) {
	key, err := getBranceKey()
	if err != nil {
		return "", err
	}

	brancaObj := branca.NewBranca(key)
	return brancaObj.EncodeToString(needEncodeString)
}

// decode
func BrancaDecode(needDecodeString string) (string, error) {
	key, err := getBranceKey()
	if err != nil {
		return "", err
	}

	brancaObj := branca.NewBranca(key)
	return brancaObj.DecodeToString(needDecodeString)
}

// Brance key
func getBranceKey() (string, error) {
	mapConfig, err := global.GetConfig("branca")
	if err != nil {
		return "", err
	}
	key, ok := mapConfig["key"].(string)
	if !ok {
		key = defKey
		global.Log.Error("获取密钥错误, 使用默认key")
	}
	return key, nil
}
