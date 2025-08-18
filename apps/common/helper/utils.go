package helper

import (
	"errors"
	"fmt"
	"regexp"
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

func CheckCellPhone(phone string) bool {
	return regexp.MustCompile(`^1[3-9]\d{8}$`).MatchString(phone)
}
