package routers

import (
	adminSysAuth "WenBeego/apps/admin_plat/controllers/auth"
	"WenBeego/apps/common/middleware"

	beego "github.com/beego/beego/v2/server/web"
	"github.com/beego/beego/v2/server/web/context"
)

// 白名单api
var WhiteApiList = []string{
	"/admin_plat/auth/login",
}

// 登录后基础api
var AuthApiList = []string{
	"/admin_plat/auth/logout",
	"/admin_plat/auth/getUserInfo",
	"/admin_plat/auth/getPermissionList",
}

func init() {
	ns := beego.NewNamespace("/admin_plat",
		beego.NSCtrlPost("/auth/login", (*adminSysAuth.AuthController).Login),
	)

	// 请求前、后处理
	ns.Filter("before", func(ctx *context.Context) {
		(new(middleware.AccessMiddleware).RouterBefore())(ctx)
		middleware.Auth(&WhiteApiList, &AuthApiList)(ctx)
	})
	// ns.Filter("after", func(ctx *context.Context) {
	// 	(new(middleware.AccessMiddleware).RouterAfter())(ctx) // 请求后处理存在bug
	// })
	beego.InsertFilter("/admin_plat/*", beego.FinishRouter, new(middleware.AccessMiddleware).RouterAfter(), beego.WithReturnOnOutput(false))

	beego.AddNamespace(ns)
}
