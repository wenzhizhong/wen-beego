package blocker

import (
	"WenBeego/apps/common/helper"
	"encoding/json"

	beecontext "github.com/beego/beego/v2/server/web/context"
)

type HeaderBaseInfo struct {
	url, host, scheme, method, token, ip, signature string
}

func getBaseInfo(ctx *beecontext.Context) HeaderBaseInfo {
	return HeaderBaseInfo{
		ip:        ctx.Request.RemoteAddr,
		url:       ctx.Request.URL.Path,
		host:      ctx.Request.Host,
		scheme:    ctx.Request.URL.Scheme,
		method:    ctx.Request.Method,
		token:     ctx.Request.Header.Get("Authorization"),
		signature: ctx.Request.Header.Get("Signature"),
	}
}
func responseStr(code int, msg string, data interface{}) string {
	res := helper.Response(code, msg, data)
	jsonString, _ := json.Marshal(res)
	return string(jsonString)
}
func setResponse(ctx *beecontext.Context, code int, msg string, data interface{}) {
	jsonString := responseStr(code, msg, data)
	ctx.ResponseWriter.ResponseWriter.WriteHeader(code)
	ctx.ResponseWriter.Write([]byte(jsonString))
}
