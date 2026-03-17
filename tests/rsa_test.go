package tests

import (
	"WenBeego/apps/common/helper"
	"fmt"
	"testing"
)

func TestRsaPrivateKey(t *testing.T) {
	privateKey, err := helper.RsaGenerate()

	fmt.Println(privateKey, err)
}

func TestRsaKeyToString(t *testing.T) {
	// 生成 RSA 密钥对
	privateKey, err := helper.RsaGenerate()
	if err != nil {
		t.Fatal(err)
	}

	// 将私钥转换为 PEM 字符串
	privateKeyStr, err := helper.RsaPrivateKeyToPem(privateKey)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("Private Key (PEM):")
	fmt.Println(privateKeyStr)

	// 获取公钥并转换为 PEM 字符串
	publicKey := &privateKey.PublicKey
	publicKeyStr, err := helper.RsaPublicKeyToPem(publicKey)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("\nPublic Key (PEM):")
	fmt.Println(publicKeyStr)

	// 测试从 PEM 字符串恢复私钥
	restoredPrivateKey, err := helper.RsaPemToPrivateKey(privateKeyStr)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("\nRestored Private Key:", restoredPrivateKey.N.String())

	// 测试从 PEM 字符串恢复公钥
	restoredPublicKey, err := helper.RsaPemToPublicKey(publicKeyStr)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("Restored Public Key:", restoredPublicKey.N.String())

	// 验证恢复的密钥是否匹配
	if restoredPrivateKey.N.Cmp(privateKey.N) != 0 {
		t.Error("Restored private key does not match original")
	}
	if restoredPublicKey.N.Cmp(publicKey.N) != 0 {
		t.Error("Restored public key does not match original")
	}
}
