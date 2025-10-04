package tests

import (
	"WenBeego/apps/common/helper"
	"fmt"
	"testing"
)

func init() {
}

func TestGetDate(t *testing.T) {
	fmt.Println(helper.GetTimeString())
	fmt.Println(helper.GetDateStamp())
	fmt.Println(helper.GetTimestamp())
	fmt.Println(helper.GetDateStamp("2022-01-01"))
	fmt.Println(helper.GetDateStamp("2022-01-02 00:00:01"))
	fmt.Println(helper.GetTimestamp("2022-01-03"))
	fmt.Println(helper.GetTimestamp("2022-01-04 00:00:01"))
}
