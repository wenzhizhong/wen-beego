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

func TestJwsSign2(t *testing.T) {
	pubKeyStr := `-----BEGIN PUBLIC KEY-----
MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAp++X3dpdlgjNIaXI0CVU
Fa+1N7qKP9u7aGeLiZRlPpp8GN3TUumcoPNNdPJH4Lc8lrHAkMdpoqRMpNCoEeTr
tNIJIT9dq1ebMjrElouXEK7/8qwA34CAQXO5ufG9BISPTsTAae0AL1q7ObB4edkq
6h5mH5EtwcIkXsg7JKO/e6H0zMQl9jk8jn7CHRb/NBK11Nl5ZgN2ZUTmsNzcgC59
/iAxh89mh+DUFiMZc00FGloQbBbbyyQS5KsuBZdu3qvgjeIvcPgFTIjLJbSkfVgs
StA7gJoBR/IHURCa1cAaOh1tJMglamKuB/ZhIpyYXn56bTtb+e+P/8lI6jNOdxby
EwIDAQAB
-----END PUBLIC KEY-----`

	// 	priKeyStr := `-----BEGIN PRIVATE KEY-----
	// MIIEvAIBADANBgkqhkiG9w0BAQEFAASCBKYwggSiAgEAAoIBAQCn75fd2l2WCM0h
	// pcjQJVQVr7U3uoo/27toZ4uJlGU+mnwY3dNS6Zyg80108kfgtzyWscCQx2mipEyk
	// 0KgR5Ou00gkhP12rV5syOsSWi5cQrv/yrADfgIBBc7m58b0EhI9OxMBp7QAvWrs5
	// sHh52SrqHmYfkS3BwiReyDsko797ofTMxCX2OTyOfsIdFv80ErXU2XlmA3ZlROaw
	// 3NyALn3+IDGHz2aH4NQWIxlzTQUaWhBsFtvLJBLkqy4Fl27eq+CN4i9w+AVMiMsl
	// tKR9WCxK0DuAmgFH8gdREJrVwBo6HW0kyCVqYq4H9mEinJhefnptO1v574//yUjq
	// M053FvITAgMBAAECggEAEWzWwmIwQsCgCaU/nHq7FP2kh+dMWrWXO7fpbnmeNcGR
	// gy0FNT215Fd1NTTB3jXKd3fIhjKYKemdZZ6cSdsJ1WXTz7DZAc3hyx5LnAlY6pYV
	// 9bKrloVUCYvop2bHbkq+P0tCtzBXIWgV3uBfkoPE+az00MK0z5La85YYLkRAD/kC
	// JLxX/H+voW3HXkbPnRar/VDS+vkZ2cSxd/0o4g2uIh++fAE8539wFlY7nuGAcv22
	// QacWHSmFbSWDYGSHqVVYQ2ZviYsb9D26uywbcnn76Bz8zoZ2T3gaZlXCcu8zIUVu
	// Bj7lwpf86p2NeRyfMRiKkqccOzBcHzc8w9Bg0nJy+QKBgQDCYB2ctDeKfoJBQ3XA
	// j5Fi5byN0w7ASTNM5//WuJXyg0qW3mmZmJPNkvWTMubgqHuf+zX2ctOdGzVMM0qF
	// HCs400QSN7iWDZZryaJcIQZPXNvMVi4ybCZOBizuG8kyTNKT0Yh4rf0O5RqI9ALt
	// wFupWQYRwx9rneLWKxzWDCOlOwKBgQDdLZf/E81TLWQx12IXVmj5/T5Z9DFcEj9r
	// cQ+8mzvaX+puUVdcQd1jjpqbDy9CvNdld0WO0YrVWGlR6pkoNRJy4iawvz/+Oatd
	// 3SDihlEO2dSmv+9oAw0p+3gFa3EL4uf0NeGkCKIzPQr2J5tn1OCUCOHgb5+5SFgA
	// k4nRvWQ5CQKBgF8rYFRRiMAuoOgDd6wIn06k3WUzaY2MSanmDcW8Ku7KicLEsz+Z
	// DQUiZ0rjKVfmJmF2Rj2ciy/pGndsxZfW6vKvviyNS7tse7Haz7v8D3LcLGIn8AaQ
	// HVEmhOkwgZo3MwNdHEy6I6UfV5amoqh1/ms0Q1x/BOtUKrRh94R1/R3xAoGAGXCP
	// FQXADhsgdSMi4zBLLsXUECCoNMDcjo0YlEb+oWV632l3tOLWhgb2/XLHqtNxqvgH
	// BiBP6a4bnxJuv1MrZg9hB99XivQzI761c5ijZiPj87IL5VjEgNmtumHbRNS6fTpd
	// U7KyhVY2Fo4Dr/OqSRykbl4obvVFOfu+VGOGTZECgYA/hTIYlkcHH4xCUx7MPQYI
	// eKacvDH+tnbmvahUAO6HDpjWSx8La6ypfKEq9ribwOO4ze2Ruv1NKqkE1q+E4UxQ
	// 6KPg00wLuo6xH0R5e5fTWYg9hvfeCl25HH9Rzz++Ru6sxfucwdH8oxuvfFR5pvdb
	// xG4at7tj4+eaq46CHLF+5A==
	// -----END PRIVATE KEY-----`

	payload := []byte("authCode1111authCodeIdpasswordG72+shD3^6phone15912345678")

	pubKey, err := helper.RsaPemToPublicKey(pubKeyStr)
	if err != nil {
		fmt.Println("pubKey err = ", err)
		return
	}

	// priKey, err := helper.RsaPemToPrivateKey(priKeyStr)
	// if err != nil {
	// 	fmt.Println("priKey err = ", err)
	// 	return
	// }
	// signerHeader := map[string]string{
	// 	"alg": string(jose.RS256),
	// }
	// sign1, err1 := goJose.JwsSign(payload, priKey, jose.RS256, true, signerHeader)
	// sign2, err2 := goJose.JwsSign(payload, priKey, jose.RS256, false, signerHeader)
	// fmt.Println("sign1= ", sign1)
	// fmt.Println("err1= ", err1)

	// fmt.Println("sign2= ", sign2)
	// fmt.Println("err2= ", err2)

	signature := "eyJhbGciOiJSUzUxMiJ9.YXV0aENvZGUxMTExYXV0aENvZGVJZHBhc3N3b3JkRzcyK3NoRDNeNnBob25lMTU5MTIzNDU2Nzg.WjsykfK-NinIgA34_rx2F6mR4rUP3xv3kVKZt57UKGs8EgY_ijxsqiiiHBHNaf0fHTk7fp19X9NkGDVeOsHRBu6jQBRAif8J9aOjpa5GJWHeOBxp7UI7gLHCjKWVRiRxpSwug3w-lGedhqC6RdfMe6ZTENFK1EvT51Z5t5FJZrKmKUIUG4JLdAJNHygnUF_Xm5GIyCsFKN_Qfti8HWUfMRJ2Swg1tQLt15Bl9KxDmoP89rgmOIsO02PhHvh4qhoufwRl9qsMUx6MglKLlJrn8w4pP197paArGMc9nNAsy6QBc4l3ec4KGSvM93LraA89mumOyR-ImHaTH_r1Y70cag"

	err = goJose.JwsSignVerify(payload, signature, pubKey, jose.RS512)
	fmt.Printf("err=%#v \n", err)

}
