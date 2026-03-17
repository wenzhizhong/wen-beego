package helper

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
)

// GenerateRsa：RSA私钥
func RsaGenerate() (*rsa.PrivateKey, error) {
	return rsa.GenerateKey(rand.Reader, 2048)
}

// RsaPrivateKeyToPem：将 RSA 私钥转换为 PEM 字符串
func RsaPrivateKeyToPem(privateKey *rsa.PrivateKey) (string, error) {
	if privateKey == nil {
		return "", errors.New("private key is nil")
	}

	privateKeyBytes := x509.MarshalPKCS1PrivateKey(privateKey) // 将私钥转换为 DER 格式

	privateKeyPem := &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: privateKeyBytes,
	}
	return string(pem.EncodeToMemory(privateKeyPem)), nil // 编码为 PEM 字符串
}

// RsaPublicKeyToPem：将 RSA 公钥转换为 PEM 字符串
func RsaPublicKeyToPem(publicKey *rsa.PublicKey) (string, error) {
	if publicKey == nil {
		return "", errors.New("public key is nil")
	}

	publicKeyBytes, err := x509.MarshalPKIXPublicKey(publicKey) // 将公钥转换为 DER 格式
	if err != nil {
		return "", err
	}

	publicKeyPem := &pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: publicKeyBytes,
	}
	return string(pem.EncodeToMemory(publicKeyPem)), nil
}

// RsaPemToPrivateKey：将 PEM 字符串转换回 RSA 私钥
func RsaPemToPrivateKey(pemString string) (*rsa.PrivateKey, error) {
	block, _ := pem.Decode([]byte(pemString))
	if block == nil {
		return nil, errors.New("failed to decode PEM block")
	}

	return x509.ParsePKCS1PrivateKey(block.Bytes) // 解析 PKCS#1 私钥
}

// RsaPemToPublicKey：将 PEM 字符串转换回 RSA 公钥
func RsaPemToPublicKey(pemString string) (*rsa.PublicKey, error) {
	block, _ := pem.Decode([]byte(pemString))
	if block == nil {
		return nil, errors.New("failed to decode PEM block")
	}

	publicKey, err := x509.ParsePKIXPublicKey(block.Bytes) // 解析 PKIX 公钥
	if err != nil {
		return nil, err
	}

	rsaPublicKey, ok := publicKey.(*rsa.PublicKey)
	if !ok {
		return nil, errors.New("not an RSA public key")
	}

	return rsaPublicKey, nil
}
