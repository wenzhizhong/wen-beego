package tests

import (
	"WenBeego/apps/common/helper"
	"WenBeego/apps/common/helper/goJose"
	"fmt"
	"testing"

	"github.com/go-jose/go-jose/v4"
)

func TestJweCrypt(t *testing.T) {
	privateKey, err := helper.RsaGenerate()
	privateKey2, _ := helper.RsaGenerate()
	if err != nil {
		fmt.Println(err)
		return
	}
	payload := []byte("hello world")

	jweString, err := goJose.JweEncrypt(payload, privateKey.Public(), jose.A128GCM, jose.RSA_OAEP, jose.NONE, nil)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("jweString= %s \n", jweString)

	parseJweSecret := []byte("")
	parseJweSecret, err = goJose.JweDecrypt(jweString, privateKey, jose.A128GCM, jose.RSA_OAEP)
	fmt.Printf("正常解密 err=%#v \n", err)
	fmt.Printf("正常解密 parseJweSecret=%#v \n", string(parseJweSecret))

	parseJweSecret, err = goJose.JweDecrypt(jweString, privateKey2, jose.A128GCM, jose.RSA_OAEP)
	fmt.Printf("错误解密1 err=%#v \n", err)
	fmt.Printf("错误解密1 parseJweSecret=%#v \n", string(parseJweSecret))

	parseJweSecret, err = goJose.JweDecrypt(jweString[2:], privateKey, jose.A128GCM, jose.RSA_OAEP)
	fmt.Printf("错误解密2 err=%#v \n", err)
	fmt.Printf("错误解密2 parseJweSecret=%#v \n", string(parseJweSecret))

	fmt.Println("\n====================")
	jweString, err = goJose.JweEncrypt(payload, privateKey.Public(), jose.A256GCM, jose.RSA_OAEP, jose.NONE, nil)
	parseJweSecret, err = goJose.JweDecrypt(jweString, privateKey, jose.A256GCM, jose.RSA_OAEP)
	fmt.Println("\n===== A256GCM jweString =====\n " + jweString)
	fmt.Println("===== A256GCM parseJweSecret =====\n" + string(parseJweSecret))

	fmt.Println("\n===== end =====")

}
