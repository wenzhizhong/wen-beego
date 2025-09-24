package routers

import (
	adminAuth "WenBeego/apps/admin_plat/controllers/auth"
	adminMenu "WenBeego/apps/admin_plat/controllers/menu"
	adminPlat "WenBeego/apps/admin_plat/controllers/plat"
	"WenBeego/apps/common/middleware"

	beego "github.com/beego/beego/v2/server/web"
	"github.com/beego/beego/v2/server/web/context"
)

// 白名单api
var WhiteApiList = []string{
	"/admin_plat/auth/login",
	"/admin_plat/auth/register",
	"/admin_plat/auth/forget-password",
	"/admin_plat/auth/get-captcha",
	"/admin_plat/auth/refresh-token",
}

// 登录后基础api
var AuthApiList = []string{
	"/admin_plat/auth/logout",
	"/admin_plat/auth/getUserInfo",
	"/admin_plat/auth/getPermissionList",
	"/admin_plat/plat/change-unit",
	"/admin_plat/plat/get-user-unit",
	"/admin_plat/menu/get-async-routes",
}

func init() {
	ns := beego.NewNamespace("/admin_plat",
		beego.NSCtrlPost("/auth/login", (*adminAuth.AuthController).Login),
		beego.NSCtrlGet("/auth/get-captcha", (*adminAuth.AuthController).GetCatpcha),
		beego.NSCtrlPost("/auth/refresh-token", (*adminAuth.AuthController).RefreshToken),
		beego.NSCtrlPost("/plat/change-unit", (*adminPlat.PlatController).ChangeUnit),
		beego.NSCtrlGet("/plat/get-user-unit-list", (*adminPlat.PlatController).GetUserUnitList),
		beego.NSCtrlGet("/menu/get-async-routes", (*adminMenu.MenuController).GetAsyncRoutes),
	)

	// 请求前、后处理
	ns.Filter("before", func(ctx *context.Context) {
		(new(middleware.AccessMiddleware).RouterBefore())(ctx)
		middleware.AuthAdmin(&WhiteApiList, &AuthApiList)(ctx)
	})
	// ns.Filter("after", func(ctx *context.Context) {
	// 	(new(middleware.AccessMiddleware).RouterAfter())(ctx) // 请求后处理存在bug
	// })
	beego.InsertFilter("/admin_plat/*", beego.FinishRouter, new(middleware.AccessMiddleware).RouterAfter(), beego.WithReturnOnOutput(false))

	beego.AddNamespace(ns)
}
