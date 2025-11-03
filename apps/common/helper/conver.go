package helper

import (
	"errors"
	"fmt"
	"strconv"
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

// interface{} 转 int64
func Interface2Int64(value interface{}) (int64, error) {
	switch v := value.(type) {
	case int:
		return int64(v), nil
	case int64:
		return v, nil
	case int32:
		return int64(v), nil
	case int16:
		return int64(v), nil
	case int8:
		return int64(v), nil
	case uint:
		return int64(v), nil
	case uint64:
		if v > 1<<63-1 {
			return 0, fmt.Errorf("uint64值过大: %d", v)
		}
		return int64(v), nil
	case uint32:
		return int64(v), nil
	case uint16:
		return int64(v), nil
	case uint8:
		return int64(v), nil
	case float64:
		if v < -1<<63 || v >= 1<<63 {
			return 0, fmt.Errorf("float64值超出范围: %f", v)
		}
		return int64(v), nil
	case float32:
		if v < -1<<63 || v >= 1<<63 {
			return 0, fmt.Errorf("float32值超出范围: %f", v)
		}
		return int64(v), nil
	case string:
		return strconv.ParseInt(v, 10, 64)
	default:
		return 0, fmt.Errorf("不支持的类型: %T", value)
	}
}

func String2Int64(value string) (int64, error) {
	return strconv.ParseInt(value, 10, 64)
}

func String2Int(value string) (int, error) {
	return strconv.Atoi(value)
}

func Int2String(value int) string {
	return strconv.Itoa(value)
}

func Int642String(value int64) string {
	return strconv.FormatInt(value, 10)
}

func Float642String(value float64) string {
	return strconv.FormatFloat(value, 'f', -1, 64)
}
