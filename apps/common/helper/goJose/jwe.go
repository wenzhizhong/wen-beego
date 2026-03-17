package goJose

import (
	"crypto/rsa"
	"errors"

	"github.com/go-jose/go-jose/v4"
)

func newEncrypter(publicKey interface{}, enc jose.ContentEncryption, alg jose.KeyAlgorithm, compressionAlgorithm jose.CompressionAlgorithm, header map[string]string) (jose.Encrypter, error) {

	rcpt := jose.Recipient{
		Algorithm: alg,
		Key:       publicKey,
	}
	opts := jose.EncrypterOptions{}

	var extraHeaders = make(map[jose.HeaderKey]interface{})
	if len(header) > 0 {
		for k, v := range header {
			extraHeaders[jose.HeaderKey(k)] = v
		}
		opts.ExtraHeaders = extraHeaders
	}
	if compressionAlgorithm != "" {
		opts.Compression = compressionAlgorithm
	}
	return jose.NewEncrypter(enc, rcpt, &opts)
}

/**
 * JWE 加密
 */
func JweEncrypt(payload []byte, publicKey interface{}, enc jose.ContentEncryption, alg jose.KeyAlgorithm, compressionAlgorithm jose.CompressionAlgorithm, header map[string]string) (string, error) {
	if len(payload) == 0 {
		return "", errors.New("payload is empty")
	}

	newEncrypter, err := newEncrypter(publicKey, enc, alg, compressionAlgorithm, header)
	if err != nil {
		return "", err
	}
	JSONWebEncryptionObj, err := newEncrypter.Encrypt(payload)
	if err != nil {
		return "", err
	}
	return JSONWebEncryptionObj.CompactSerialize()
}

/**
 * JWE 解密
 *
 */

func JweDecrypt(jweString string, privateKey *rsa.PrivateKey, enc jose.ContentEncryption, alg jose.KeyAlgorithm) ([]byte, error) {
	if len(jweString) == 0 {
		return nil, errors.New("jweString is empty")
	}

	contentEncryption := []jose.ContentEncryption{enc}
	keyEncryptionAlgorithms := []jose.KeyAlgorithm{alg}
	JSONWebEncryptionObj, err := jose.ParseEncrypted(jweString, keyEncryptionAlgorithms, contentEncryption)
	if err != nil {
		return nil, err
	}
	return JSONWebEncryptionObj.Decrypt(privateKey)
}
