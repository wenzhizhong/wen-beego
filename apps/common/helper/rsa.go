package helper

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"strings"
)

// ============================================================
// 一、密钥生成
// ============================================================

// RsaGenerate 生成 RSA 私钥（2048 位）
func RsaGenerate() (*rsa.PrivateKey, error) {
	return rsa.GenerateKey(rand.Reader, 2048)
}

// RsaGenerateWithBits 生成指定密钥位数的 RSA 私钥
// 安全建议：bits >= 2048；高安全场景可用 4096
func RsaGenerateWithBits(bits int) (*rsa.PrivateKey, error) {
	if bits < 1024 {
		return nil, errors.New("rsa: bits too small, must be >= 1024")
	}
	return rsa.GenerateKey(rand.Reader, bits)
}

// ============================================================
// 二、PEM 序列化（对象 -> PEM 字符串）
// ============================================================

// RsaPrivateKeyToPem 私钥 -> PKCS#8 PEM（"PRIVATE KEY"）
func RsaPrivateKeyToPem(privateKey *rsa.PrivateKey) (string, error) {
	if privateKey == nil {
		return "", errors.New("private key is nil")
	}

	privateKeyBytes, err := x509.MarshalPKCS8PrivateKey(privateKey)
	if err != nil {
		return "", err
	}

	privateKeyPem := &pem.Block{
		Type:  "PRIVATE KEY",
		Bytes: privateKeyBytes,
	}
	return string(pem.EncodeToMemory(privateKeyPem)), nil
}

// RsaPrivateKeyToPemPKCS1 私钥 -> PKCS#1 PEM（"RSA PRIVATE KEY"）
// 用于对接 OpenSSL 传统格式 / 部分老系统
func RsaPrivateKeyToPemPKCS1(privateKey *rsa.PrivateKey) (string, error) {
	if privateKey == nil {
		return "", errors.New("private key is nil")
	}

	privateKeyBytes := x509.MarshalPKCS1PrivateKey(privateKey)
	privateKeyPem := &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: privateKeyBytes,
	}
	return string(pem.EncodeToMemory(privateKeyPem)), nil
}

// RsaPublicKeyToPem 公钥 -> PKIX PEM（"PUBLIC KEY"）
func RsaPublicKeyToPem(publicKey *rsa.PublicKey) (string, error) {
	if publicKey == nil {
		return "", errors.New("public key is nil")
	}

	publicKeyBytes, err := x509.MarshalPKIXPublicKey(publicKey)
	if err != nil {
		return "", err
	}

	publicKeyPem := &pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: publicKeyBytes,
	}
	return string(pem.EncodeToMemory(publicKeyPem)), nil
}

// RsaPublicKeyToPemPKCS1 公钥 -> PKCS#1 PEM（"RSA PUBLIC KEY"）
// 用于对接 OpenSSL 传统格式 / 部分老系统
func RsaPublicKeyToPemPKCS1(publicKey *rsa.PublicKey) (string, error) {
	if publicKey == nil {
		return "", errors.New("public key is nil")
	}

	publicKeyBytes := x509.MarshalPKCS1PublicKey(publicKey)
	publicKeyPem := &pem.Block{
		Type:  "RSA PUBLIC KEY",
		Bytes: publicKeyBytes,
	}
	return string(pem.EncodeToMemory(publicKeyPem)), nil
}

// ============================================================
// 三、PEM 反序列化（PEM 字符串 -> 对象）
// ============================================================

// RsaPemToPrivateKey PEM -> *rsa.PrivateKey
// 同时支持 PKCS#8（"PRIVATE KEY"）和 PKCS#1（"RSA PRIVATE KEY"）
//
// 【注意】返回类型由原来的 any 改为 *rsa.PrivateKey，
// 老的 `key.(*rsa.PrivateKey)` 断言可以删掉。
func RsaPemToPrivateKey(pemString string) (*rsa.PrivateKey, error) {
	pemString = RsaPrivateKey(pemString)
	block, _ := pem.Decode([]byte(pemString))
	if block == nil {
		return nil, errors.New("failed to decode PEM block")
	}

	switch block.Type {
	case "RSA PRIVATE KEY":
		return x509.ParsePKCS1PrivateKey(block.Bytes)

	case "PRIVATE KEY":
		key, err := x509.ParsePKCS8PrivateKey(block.Bytes)
		if err != nil {
			return nil, err
		}
		rsaKey, ok := key.(*rsa.PrivateKey)
		if !ok {
			return nil, errors.New("not an RSA private key")
		}
		return rsaKey, nil

	default:
		// 兜底：先按 PKCS#8 试，再按 PKCS#1 试
		if key, err := x509.ParsePKCS8PrivateKey(block.Bytes); err == nil {
			if rsaKey, ok := key.(*rsa.PrivateKey); ok {
				return rsaKey, nil
			}
		}
		return x509.ParsePKCS1PrivateKey(block.Bytes)
	}
}

// RsaPemToPublicKey PEM -> *rsa.PublicKey
// 同时支持 PKIX（"PUBLIC KEY"）和 PKCS#1（"RSA PUBLIC KEY"）
func RsaPemToPublicKey(pemString string) (*rsa.PublicKey, error) {
	pemString = RsaPrivateKey(pemString)
	block, _ := pem.Decode([]byte(pemString))
	if block == nil {
		return nil, errors.New("failed to decode PEM block")
	}

	switch block.Type {
	case "PUBLIC KEY":
		pub, err := x509.ParsePKIXPublicKey(block.Bytes)
		if err != nil {
			return nil, err
		}
		rsaPub, ok := pub.(*rsa.PublicKey)
		if !ok {
			return nil, errors.New("not an RSA public key")
		}
		return rsaPub, nil

	case "RSA PUBLIC KEY":
		return x509.ParsePKCS1PublicKey(block.Bytes)

	default:
		// 兜底：先按 PKIX 试，再按 PKCS#1 试
		if pub, err := x509.ParsePKIXPublicKey(block.Bytes); err == nil {
			if rsaPub, ok := pub.(*rsa.PublicKey); ok {
				return rsaPub, nil
			}
		}
		return x509.ParsePKCS1PublicKey(block.Bytes)
	}
}

// ============================================================
// 四、加密 / 解密（OAEP + SHA-256）
// ============================================================

// RsaEncrypt 使用 RSA 公钥加密（OAEP + SHA-256，推荐方案）
//
// 单次明文长度上限 = keySize - 2*sha256.Size - 2：
//   - 2048 位密钥 ≈ 190 字节
//   - 4096 位密钥 ≈ 446 字节
//
// 超长明文请使用混合加密（RSA + AES）: helper.HybridEncrypt(......)。
func RsaEncrypt(publicKey *rsa.PublicKey, plaintext []byte) ([]byte, error) {
	if publicKey == nil {
		return nil, errors.New("public key is nil")
	}
	return rsa.EncryptOAEP(sha256.New(), rand.Reader, publicKey, plaintext, nil)
}

// RsaDecrypt 使用 RSA 私钥解密（OAEP + SHA-256）
func RsaDecrypt(privateKey *rsa.PrivateKey, ciphertext []byte) ([]byte, error) {
	if privateKey == nil {
		return nil, errors.New("private key is nil")
	}
	return rsa.DecryptOAEP(sha256.New(), rand.Reader, privateKey, ciphertext, nil)
}

// RsaMaxEncryptLen 返回 OAEP + SHA-256 下单次可加密的最大明文长度
func RsaMaxEncryptLen(publicKey *rsa.PublicKey) int {
	if publicKey == nil {
		return 0
	}
	return publicKey.Size() - 2*sha256.Size - 2
}

// ============================================================
// 五、签名 / 验签（PSS + SHA-256）
// ============================================================

// RsaSign 使用私钥对消息签名（PSS + SHA-256，推荐方案）
// 入参为原始消息，函数内部自动做 SHA-256 哈希
func RsaSign(privateKey *rsa.PrivateKey, message []byte) ([]byte, error) {
	if privateKey == nil {
		return nil, errors.New("private key is nil")
	}
	hashed := sha256.Sum256(message)
	return rsa.SignPSS(rand.Reader, privateKey, crypto.SHA256, hashed[:], nil)
}

// RsaVerify 使用公钥验签（PSS + SHA-256）
// 返回 nil 表示验签通过；非 nil 表示失败
func RsaVerify(publicKey *rsa.PublicKey, message, signature []byte) error {
	if publicKey == nil {
		return errors.New("public key is nil")
	}
	hashed := sha256.Sum256(message)
	return rsa.VerifyPSS(publicKey, crypto.SHA256, hashed[:], signature, nil)
}

// RsaSignPKCS1v15 使用 PKCS#1 v1.5 + SHA-256 签名（仅用于对接旧系统）
func RsaSignPKCS1v15(privateKey *rsa.PrivateKey, message []byte) ([]byte, error) {
	if privateKey == nil {
		return nil, errors.New("private key is nil")
	}
	hashed := sha256.Sum256(message)
	return rsa.SignPKCS1v15(rand.Reader, privateKey, crypto.SHA256, hashed[:])
}

// RsaVerifyPKCS1v15 使用 PKCS#1 v1.5 + SHA-256 验签（仅用于对接旧系统）
func RsaVerifyPKCS1v15(publicKey *rsa.PublicKey, message, signature []byte) error {
	if publicKey == nil {
		return errors.New("public key is nil")
	}
	hashed := sha256.Sum256(message)
	return rsa.VerifyPKCS1v15(publicKey, crypto.SHA256, hashed[:], signature)
}

// ============================================================
// 六、PEM 字符串清理
// ============================================================

// RsaPrivateKey 清理 PEM 字符串中的 \r、\t 等
func RsaPrivateKey(pemString string) string {
	pemString = strings.ReplaceAll(pemString, "\r", "")
	pemString = strings.ReplaceAll(pemString, "\t", "")
	// pemString = strings.ReplaceAll(pemString, "\n", "\\n")
	return pemString

}
