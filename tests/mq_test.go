package tests

import (
	"WenBeego/apps/common/middleware/mq"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/RichardKnop/machinery/v1/tasks"
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
	}

	asyncResult, err := (&mq.MqClient{}).SendTask("test", options)
	fmt.Println("\n", asyncResult, err)
}
