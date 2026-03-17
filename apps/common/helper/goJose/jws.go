package goJose

// go-jose 封装
import (
	"errors"

	"github.com/go-jose/go-jose/v4"
)

/**
 * 实例化jose Signer
 *
 * privateKey  ed25519.PrivateKey | *rsa.PrivateKey | *ecdsa.PrivateKey | []byte | JSONWebKey | *JSONWebKey | OpaqueSigner
 * alg  jose.RS256| jose.RS384| jose.RS512| jose.PS256| jose.PS384| jose.PS512
 * */
func newSigner(privateKey interface{}, alg jose.SignatureAlgorithm, embedJWK bool, signerHeader map[string]string) (jose.Signer, error) {
	sig := jose.SigningKey{
		Algorithm: alg,
		Key:       privateKey,
	}
	opts := jose.SignerOptions{
		NonceSource:  &NonceSourceStruct{},
		EmbedJWK:     embedJWK,
		ExtraHeaders: nil,
	}

	extraHeaders := make(map[jose.HeaderKey]interface{})
	if len(signerHeader) > 0 {
		for k, v := range signerHeader {
			extraHeaders[jose.HeaderKey(k)] = v
		}
	}
	opts.ExtraHeaders = extraHeaders

	return jose.NewSigner(sig, &opts)
}

/**
 * jose 签名
 *
 * payload
 * privateKey  ed25519.PrivateKey | *rsa.PrivateKey | *ecdsa.PrivateKey | []byte | JSONWebKey | *JSONWebKey | OpaqueSigner
 * alg  jose.RS256| jose.RS384| jose.RS512| jose.PS256| jose.PS384| jose.PS512
 * */
func JwsSign(payload []byte, privateKey interface{}, alg jose.SignatureAlgorithm, embedJWK bool, signerHeader map[string]string) (string, error) {
	if len(payload) == 0 {
		return "", errors.New("payload is empty")
	}
	signer, err := newSigner(privateKey, alg, embedJWK, signerHeader)
	if err != nil {
		return "", err
	}

	JsonWebSignature, err := signer.Sign(payload)
	if err != nil {
		return "", err
	}
	return JsonWebSignature.DetachedCompactSerialize()
}

/**
 * jose 校验签名
 *
 * */
func JwsSignVerify(payload []byte, signature string, publicKey interface{}, alg jose.SignatureAlgorithm) error {
	JsonWebSignature, err := jose.ParseDetached(signature, payload, []jose.SignatureAlgorithm{alg})
	if err != nil {
		return err
	}
	return JsonWebSignature.DetachedVerify(payload, publicKey)
}
