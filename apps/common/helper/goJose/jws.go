package goJose

// go-jose 封装
import (
	"errors"
	"fmt"
	"sort"
	"strings"

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

// SerializeForPayload 将任意 JSON 兼容数据转换为用于签名的字符串
func SerializeForPayload(v interface{}) string {
	switch val := v.(type) {
	case nil:
		return "null"
	case string:
		return val
	case bool:
		if val {
			return "true"
		}
		return "false"
	case float64, float32, int, int64, int32, int16, int8, uint, uint64, uint32, uint16, uint8:
		return fmt.Sprint(val) // 数字转为十进制字符串
	case []interface{}:
		elems := make([]string, len(val))
		for i, e := range val {
			elems[i] = SerializeForPayload(e)
		}
		return "[" + strings.Join(elems, ",") + "]"
	case map[string]interface{}:
		keys := make([]string, 0, len(val))
		for k := range val {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		var parts []string
		for _, k := range keys {
			parts = append(parts, k+SerializeForPayload(val[k]))
		}
		return strings.Join(parts, "")
	default:
		// 其他类型（如通过反射处理复杂结构），可根据实际扩展
		panic(fmt.Sprintf("unsupported type: %T", v))
	}
}
