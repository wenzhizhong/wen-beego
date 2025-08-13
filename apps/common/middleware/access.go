package middleware

import (
	"fmt"
	"net/http"

	"github.com/beego/beego/v2/core/logs"
	"github.com/beego/beego/v2/server/web"
	beecontext "github.com/beego/beego/v2/server/web/context"
	"golang.org/x/time/rate"
)

var limiter = rate.NewLimiter(10, 20)

type AccessMiddleware struct {
}

// 限制访问次数,防止dos攻击
func (m *AccessMiddleware) LimitTimes() web.FilterFunc {
	return func(ctx *beecontext.Context) {
		if !limiter.Allow() {
			fmt.Println("Too Many Requests")
			logs.Warn("Too Many Requests")
			http.Error(ctx.ResponseWriter, "Too Many Requests", http.StatusTooManyRequests)
			return
		}
	}
}

// api 请求前
func (m *AccessMiddleware) RouterBefore() web.FilterFunc {
	return func(ctx *beecontext.Context) {
		// url, host, shceme, method, token, ip := m.getBaseInfo(ctx)
		// timeStr := time.Now().Format("2006-01-02 15:04:05")
	}
}

// api 请求后
func (m *AccessMiddleware) RouterAfter() web.FilterFunc {
	return func(ctx *beecontext.Context) {
		// url, host, shceme, method, token, ip := m.getBaseInfo(ctx)
		// timeStr := time.Now().Format("2006-01-02 15:04:05")
	}
}

func (m *AccessMiddleware) getBaseInfo(ctx *beecontext.Context) (url, host, shceme, method, token, ip string) {
	url = ctx.Request.URL.Path
	host = ctx.Request.Host
	shceme = ctx.Request.URL.Scheme
	method = ctx.Request.Method
	token = ctx.Request.Header.Get("Authorization")
	ip = ctx.Request.RemoteAddr
	return
}
