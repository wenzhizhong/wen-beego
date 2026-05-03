package controllers

import "fmt"

type Test struct {
}

func (c *Test) ActionTestMsg(username, hobby string, age string) error {
	fmt.Println("username:", username, "hobby:", hobby, "age:", age, "\n==============")

	var err = fmt.Errorf("ActionTestMsg(): 测试错误 %s", username)
	return err

}
