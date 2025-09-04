package tests

import (
	"WenBeego/apps/common/helper"
	"fmt"
	"testing"
)

func init() {
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
