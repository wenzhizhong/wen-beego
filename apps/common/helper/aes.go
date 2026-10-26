package helper

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/rsa"
	"errors"
)

// HybridEncrypt 发送方使用：用接收方的 RSA 公钥加密任意长度的明文
//
// 返回：
//   - encryptedKey：AES 密钥的 RSA 密文（网络上传输这个，不是 AES 明文密钥）
//   - ciphertext  ：AES-GCM 密文（nonce 已拼在密文头部）
func HybridEncrypt(publicKey *rsa.PublicKey, plaintext []byte) (encryptedKey, ciphertext []byte, err error) {
	// 1. 本地随机生成 32 字节 AES-256 密钥，只保存在内存里
	aesKey := make([]byte, 32)
	if _, err = rand.Read(aesKey); err != nil {
		return nil, nil, err
	}

	// 2. 用接收方的 RSA 公钥加密 AES 密钥（传输的是密文，不是明文密钥）
	encryptedKey, err = RsaEncrypt(publicKey, aesKey)
	if err != nil {
		return nil, nil, err
	}

	// 3. 用 AES-GCM 加密业务数据
	block, err := aes.NewCipher(aesKey)
	if err != nil {
		return nil, nil, err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, nil, err
	}
	nonce := make([]byte, gcm.NonceSize())
	if _, err = rand.Read(nonce); err != nil {
		return nil, nil, err
	}
	ciphertext = gcm.Seal(nonce, nonce, plaintext, nil) // nonce 拼在密文前面
	return encryptedKey, ciphertext, nil
}

// HybridDecrypt 接收方使用：用自己的 RSA 私钥解密
func HybridDecrypt(privateKey *rsa.PrivateKey, encryptedKey, ciphertext []byte) ([]byte, error) {
	// 1. 用自己的 RSA 私钥解出 AES 密钥
	aesKey, err := RsaDecrypt(privateKey, encryptedKey)
	if err != nil {
		return nil, err
	}

	// 2. 用 AES-GCM 解密业务数据
	block, err := aes.NewCipher(aesKey)
	if err != nil {
		return nil, err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	nonceSize := gcm.NonceSize()
	if len(ciphertext) < nonceSize {
		return nil, errors.New("ciphertext too short")
	}
	nonce, realCiphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
	return gcm.Open(nil, nonce, realCiphertext, nil)
}
