package controllers

import (
	"WenBeego/apps/common/global"
	"fmt"
)

// 生日提醒
type BirthNotice struct {
}

func (c *BirthNotice) Notice() {
	fmt.Println("birth notice .....")
	global.Log.Info("birth notice .....")
	panic("==== panic test ===")
}
