package tests

import (
	"WenBeego/apps/common/helper"
	"fmt"
	"testing"
	"time"

	gofrsUuid "github.com/gofrs/uuid"
	googleUuid "github.com/google/uuid"
)

func init() {
}

func TestGofrsUuid(t *testing.T) {
	fmt.Println(time.Now().UnixNano())
	for i := 0; i < 100; i++ {
		uuid, err := gofrsUuid.NewV7()
		if err != nil {
			fmt.Println(err.Error())
			continue
		}
		fmt.Println(uuid.String())
	}
	fmt.Println(time.Now().UnixNano())
}

func TestGoogleUuid(t *testing.T) {
	fmt.Println(time.Now().UnixNano())
	googleUuid.EnableRandPool()
	for i := 0; i < 100; i++ {
		uuid, err := googleUuid.NewV7()
		if err != nil {
			fmt.Println(err.Error())
			continue
		}
		fmt.Println(uuid.String())

	}
	fmt.Println(time.Now().UnixNano())
}

func TestStrIsUuid(t *testing.T) {
	res := helper.StringIsUuid("")
	res1 := helper.StringIsUuid("01a0533f-c0c0-77ba")
	res2 := helper.StringIsUuid("01a0533f-c0c0-77ba-a3af-62db7e33152e")
	fmt.Println(res, res1, res2)
	fmt.Println("==============")
}
