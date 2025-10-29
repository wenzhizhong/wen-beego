package tests

import (
	"WenBeego/apps/common/helper"
	"fmt"
	"testing"
)

func init() {
}

func TestGetDate(t *testing.T) {
	fmt.Println("helper.GetTimeString()\t\t\t\t", helper.GetTimeString())

	fmt.Println("helper.GetTimeString(\"2006/01/02\")\t\t", helper.GetTimeString("2006/01/02"))

	fmt.Println("helper.GetTimeString(\"2006/01/02 00:00:00\")\t", helper.GetTimeString("2006/01/02 00:00:00"))

	fmt.Println("helper.GetTimeString(\"2006/01/02 00:00:05\")\t", helper.GetTimeString("2006/01/02 00:00:05"))

	fmt.Println("helper.GetDateStamp()\t\t\t\t", helper.GetDateStamp(), "\t\t\t helper.TimestampToTime()", helper.TimestampToTime(helper.GetDateStamp()))

	fmt.Println("helper.GetTimestamp()\t\t\t\t", helper.GetTimestamp(), "\t\t\t helper.TimestampToTime()", helper.TimestampToTime(helper.GetTimestamp()))

	fmt.Println("helper.GetDateStamp(\"2025-10-29\")\t\t", helper.GetDateStamp("2025-10-29"), "\t\t\t helper.TimestampToTime()", helper.TimestampToTime(helper.GetDateStamp("2025-10-29")))

	fmt.Println("helper.GetDateStamp(\"2025-10-29 12:05:06\")\t", helper.GetDateStamp("2025-10-29 12:05:06"), "\t\t\t helper.TimestampToTime()", helper.TimestampToTime(helper.GetDateStamp("2025-10-29 12:05:06")))

	fmt.Println("helper.GetTimestamp(\"2025-10-29\")\t\t", helper.GetTimestamp("2025-10-29"), "\t\t\t helper.TimestampToTime()", helper.TimestampToTime(helper.GetTimestamp("2025-10-29")))

	fmt.Println("helper.GetTimestamp(\"2025-10-29 12:05:06\")\t", helper.GetTimestamp("2025-10-29 12:05:06"), "\t\t\t helper.TimestampToTime()", helper.TimestampToTime(helper.GetTimestamp("2025-10-29 12:05:06")))

	fmt.Println("==============end ==================")

}
