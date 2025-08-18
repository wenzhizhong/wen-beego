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

func TestValidatePassword(t *testing.T) {
	fmt.Println("z<Xj.Z7450 输出", helper.ValidatePassword("z<Xj.Z7450"))
	fmt.Println("123456 输出", helper.ValidatePassword("123456"))
	fmt.Println("abcdefgh123efgh123456 输出", helper.ValidatePassword("abcdefgh123efgh123456"))
	fmt.Println("abcdefgh 输出", helper.ValidatePassword("abcdefgh"))
	fmt.Println("123456789 输出", helper.ValidatePassword("123456789"))
	fmt.Println("aB123456789< 输出", helper.ValidatePassword("aB123456789<"))
	fmt.Println("123cDEFGHI< 输出", helper.ValidatePassword("123cDEFGHI<"))

}
