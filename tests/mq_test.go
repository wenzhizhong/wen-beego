package tests

import (
	"WenBeego/apps/common/global/constant"
	"WenBeego/apps/common/middleware/mq"
	"encoding/json"
	"fmt"
	"testing"

	"WenBeego/apps/common/thirdPkg/rewrite/RichardKnop/machinery/v1/tasks"
)

func init() {
}

func TestMqSendTask(t *testing.T) {
	hobby := []struct {
		Name string `json:"name"`
	}{
		{Name: "打篮球"},
	}
	hobbyJson, _ := json.Marshal(hobby)
	fmt.Println(hobbyJson)
	options := []tasks.Arg{
		{Name: "username", Type: "string", Value: "张三"},
		{Name: "hobby", Type: "string", Value: hobbyJson},
		// {Name: "age", Type: "int64", Value: 23}, // amqp.go:363 Task failed: Reflect task args error: 23 is not int64
		{Name: "age", Type: "string", Value: "23"},
	}

	asyncResult, err := (&mq.MqClient{}).SendTask(constant.MQ_TEST_MSG, options)
	fmt.Println("\n", asyncResult, err)
}
