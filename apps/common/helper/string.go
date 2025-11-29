package helper

import "regexp"

// 删除字符串中的空格
func DeleteSpace(str string) string {
	reg := regexp.MustCompile(`\s+`)
	str = reg.ReplaceAllString(str, "")
	return str
}
