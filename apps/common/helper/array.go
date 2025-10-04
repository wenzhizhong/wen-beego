package helper

import (
	"errors"
	"reflect"
)

// 判断字符串是否在数组中
func InArray(val string, array interface{}) (bool, error) {
	reflectType := reflect.TypeOf(array)
	if !(reflectType.Kind() == reflect.Slice || reflectType.Kind() == reflect.Array) {
		return false, errors.New("InArray(): array must be slice or array")
	}
	newArray := reflect.ValueOf(array)
	for i := 0; i < newArray.Len(); i++ {
		if val == newArray.Index(i).Interface() {
			return true, nil
		}
	}
	return false, nil
}

// 数组合并
func ArrayMerge(arr1, arr2 []string) []string {
	return append(arr1, arr2...)
}
