package tests

import (
	"WenBeego/apps/common/helper"
	"WenBeego/apps/common/helper/goJose"
	"fmt"
	"testing"

	"github.com/go-jose/go-jose/v4"
)

func TestJwsSign(t *testing.T) {
	privateKey, err := helper.RsaGenerate()
	if err != nil {
		fmt.Println(err)
		return
	}
	payload := []byte("hello world")

	signature, err := goJose.JwsSign(payload, privateKey, jose.RS512, true, nil)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("signature=%s \n", signature)

	err = goJose.JwsSignVerify(payload, signature, privateKey.Public(), jose.RS512)
	fmt.Printf("正常校验 err=%#v \n", err)

	err = goJose.JwsSignVerify(payload, signature+"a", privateKey.Public(), jose.RS512)
	fmt.Printf("错误校验1 err=%#v \n", err)

	err = goJose.JwsSignVerify(payload[2:], signature, privateKey.Public(), jose.RS512)
	fmt.Printf("错误校验2 err=%#v \n", err)

	fmt.Println("===== end =====")

}
