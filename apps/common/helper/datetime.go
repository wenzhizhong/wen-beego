package helper

import (
	"WenBeego/apps/common/global"
	"time"
)

// 获取时间戳
func GetTimestamp(datetime ...string) int64 {
	timezone, _ := global.GetConfigDiy("timezone")
	loc, _ := time.LoadLocation(timezone.(string))

	if len(datetime) > 0 {
		if len(datetime[0]) == 10 {
			datetime[0] += " 00:00:00"
		}
		t, _ := time.ParseInLocation("2006-01-02 15:04:05", datetime[0], loc)
		return t.Unix()
	}
	return time.Now().In(loc).Unix()
}

// 获取日期时间戳
func GetDateStamp(date ...string) int64 {
	timezone, _ := global.GetConfigDiy("timezone")
	loc, _ := time.LoadLocation(timezone.(string))
	if len(date) > 0 {
		if len(date[0]) > 10 {
			date[0] = date[0][:10]
		}
		t, _ := time.ParseInLocation("2006-01-02", date[0], loc)
		return t.Unix()
	}
	todayString := time.Now().In(loc).Format("2006-01-02")
	today, _ := time.ParseInLocation("2006-01-02", todayString, loc)
	return today.Unix()
}

// 获取时间
func GetTime() time.Time {
	timezone, _ := global.GetConfigDiy("timezone")
	loc, _ := time.LoadLocation(timezone.(string))
	return time.Now().In(loc)
}

// 获取时间字符串
func GetTimeString(format ...string) string {
	formatStr := "2006-01-02 15:04:05"
	if len(format) > 0 {
		formatStr = format[0]
	}
	return GetTime().Format(formatStr)
}

// 时间戳转时间
func TimestampToTime(timestamp int64, format ...string) string {
	formatStr := "2006-01-02 15:04:05"
	if len(format) > 0 {
		formatStr = format[0]
	}
	return time.Unix(timestamp, 0).Format(formatStr)
}
