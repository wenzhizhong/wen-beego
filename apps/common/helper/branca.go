package helper

import (
	"WenBeego/apps/common/global"
	"errors"

	"github.com/hako/branca"
)

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
		global.Log.Error("请配置branca key")
		return "", errors.New("请配置branca key")
	}
	return key, nil
}
