package helper

import (
	"encoding/json"
	"strings"

	beecontext "github.com/beego/beego/v2/server/web/context"
)

// 获取请求的token
func GetReqToken(ctx beecontext.Context) string {
	token := ctx.Request.Header.Get("Authorization")
	token = strings.TrimPrefix(token, "Bearer ")
	return token
}

// 获取请求接口
func GetReqUrl(ctx beecontext.Context) string {
	return ctx.Request.URL.Path
}

// 获取请求方法
func GetReqMethod(ctx beecontext.Context) string {
	return ctx.Request.Method
}

// 获取请求协议
func GetReqScheme(ctx beecontext.Context) string {
	return ctx.Request.URL.Scheme
}

// 获取请求主机
func GetReqHost(ctx beecontext.Context) string {
	return ctx.Request.Host
}

func GetReqIp(ctx beecontext.Context) string {
	return ctx.Request.RemoteAddr
}

func GetReqBody[T any](ctx *beecontext.Context) (T, error) {
	var data T
	reqBody := ctx.Input.RequestBody
	if len(reqBody) > 0 {
		err := json.Unmarshal(reqBody, &data)
		if err != nil {
			return data, err
		}
	}
	return data, nil
}
