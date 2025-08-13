package middleware

import (
	"net/http"

	"github.com/beego/beego/v2/server/web"
	beecontext "github.com/beego/beego/v2/server/web/context"
)

type AuthMiddleware struct {
}

func (m *AuthMiddleware) Next() web.FilterFunc {
	return func(ctx *beecontext.Context) {
		token := ctx.Request.Header.Get("Authorization")
		if token != "valid_token" {
			ctx.ResponseWriter.WriteHeader(http.StatusUnauthorized)
			ctx.ResponseWriter.Write([]byte("未授权"))
			return
		}
	}
}
