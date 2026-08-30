package helper

import "regexp"

// 删除字符串中的空格
func DeleteSpace(str string) string {
	reg := regexp.MustCompile(`\s+`)
	str = reg.ReplaceAllString(str, "")
	return str
}

// 按字符截取字符串，支持中文
func SubStr(s string, start, length int) string {
	// 1. 字符串转 rune 切片（每个元素=1个字符）
	r := []rune(s)

	// 边界校验（防止越界）
	if start < 0 {
		start = 0
	}
	if start >= len(r) {
		return ""
	}
	end := start + length
	if end > len(r) {
		end = len(r)
	}

	// 2. 按字符截取，再转回字符串
	return string(r[start:end])
}
