package tests

import (
	"WenBeego/apps/common/helper"
	"fmt"
	"reflect"
	"testing"
)

func init() {
	fmt.Println("init .... ")
}

func TestGetRandomPassword(t *testing.T) {
	pwd, err := helper.GetRandomPassword(10)
	fmt.Println("pwd=", pwd)
	fmt.Println("err=", err)
}

func TestCheckPasswordRule(t *testing.T) {
	fmt.Println("z<Xj.Z7450 输出", helper.CheckPasswordRule("z<Xj.Z7450"))
	fmt.Println("123456 输出", helper.CheckPasswordRule("123456"))
	fmt.Println("abcdefgh123efgh123456 输出", helper.CheckPasswordRule("abcdefgh123efgh123456"))
	fmt.Println("abcdefgh 输出", helper.CheckPasswordRule("abcdefgh"))
	fmt.Println("123456789 输出", helper.CheckPasswordRule("123456789"))
	fmt.Println("aB123456789< 输出", helper.CheckPasswordRule("aB123456789<"))
	fmt.Println("123cDEFGHI< 输出", helper.CheckPasswordRule("123cDEFGHI<"))

}

func TestGetCapture(t *testing.T) {
	id, b64s, answer, err := helper.GetCaptcha("string")
	fmt.Printf("id:%s\nb64s:%s\nanswer:%s\nerr:%v\n", id, b64s, answer, err)
}

func TestVeriryCapture(t *testing.T) {
	fmt.Println(helper.VerifyCaptcha("string", "htDcENyykTRc1NH6x23y", "7650"))
}

func TestIsPhone(t *testing.T) {
	phomeNo := "15912345678"
	fmt.Printf("is phone No. %s:\n", phomeNo)
	fmt.Println(helper.IsCellPhone(phomeNo))
}

func TestGetMapKeys(t *testing.T) {
	m := map[string]string{"a": "1", "b": "2", "c": "3"}
	fmt.Println(m)
	fmt.Println(helper.GetMapKeys(m))

	m1 := map[int]string{1: "1", 2: "2", 3: "3"}
	fmt.Println(m1)
	fmt.Println(helper.GetMapKeys(m1))

	m2 := map[string]interface{}{"a": "1", "b": "2", "c": "3"}
	fmt.Println(m2)
	fmt.Println(helper.GetMapKeys(m2))

	m3 := map[string]map[string]interface{}{"a": {"aa": "1", "bb": "2", "cc": "3"}, "b": {"aa": "1", "bb": "2", "cc": "3"}, "c": {"aa": "1", "bb": "2", "cc": "3"}}
	fmt.Println(m3)
	fmt.Println(helper.GetMapKeys(m3))
}

func TestRedis(t *testing.T) {
	helper.RedisPut("test", "test", 300)
	helper.RedisPut("test1", 1, 300)
	helper.RedisPut("test2", struct {
		Test string `json:"test"`
	}{Test: "1"}, 300)

	value1, err1 := helper.RedisGet("test")
	value2, err2 := helper.RedisGet("test1")
	value3, err3 := helper.RedisGet("test2")

	fmt.Println(err1, err2, err3)
	fmt.Println(value1, value2, value3)

}

func Test(t *testing.T) {
	data1 := helper.Ternary(true, "data_true", "data_false")
	fmt.Println("data1=", data1+"1")
	fmt.Println("data1 type=", reflect.TypeOf(data1))

	var1 := "var_1"
	var2 := "var_2"
	data2 := helper.TernaryPtr(true, &var1, &var2)
	fmt.Println("data2=", data2)
	fmt.Println("data2=", *data2)
	fmt.Println("data2 type=", reflect.TypeOf(*data2))
}
