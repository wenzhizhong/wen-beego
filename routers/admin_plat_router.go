package routers

import (
	adminAuth "WenBeego/apps/admin_plat/controllers/auth"
	adminSystem "WenBeego/apps/admin_plat/controllers/system"
	adminUpload "WenBeego/apps/admin_plat/controllers/upload"
	"WenBeego/apps/common/middleware"

	beego "github.com/beego/beego/v2/server/web"
	"github.com/beego/beego/v2/server/web/context"
)

// 白名单api
var platWhiteApiList = []string{
	"/admin_plat/auth/login",
	"/admin_plat/auth/register",
	"/admin_plat/auth/forget-password",
	"/admin_plat/auth/get-captcha",
	"/admin_plat/auth/refresh-token",
}

// 登录后基础api
var platAuthApiList = []string{
	"/admin_plat/auth/logout",
	"/admin_plat/auth-params/model-params",
	"/admin_plat/auth-plat/change-unit",
	"/admin_plat/auth-plat/get-user-unit-list",
	"/admin_plat/auth-menu/get-async-routes",
	"/admin_plat/upload/upload",
	"/admin_plat/upload/vue-slice-upload",
}

func init() {
	ns := beego.NewNamespace("/admin_plat",
		// auth
		beego.NSCtrlPost("/auth/login", (*adminAuth.AuthController).Login),
		beego.NSCtrlGet("/auth/get-captcha", (*adminAuth.AuthController).GetCatpcha),
		beego.NSCtrlPost("/auth/refresh-token", (*adminAuth.AuthController).RefreshToken),
		beego.NSCtrlGet("/auth-params/model-params", (*adminAuth.ParamsController).GetModelParams),
		beego.NSCtrlPost("/auth-plat/change-unit", (*adminAuth.PlatController).ChangeUnit),
		beego.NSCtrlGet("/auth-plat/get-user-unit-list", (*adminAuth.PlatController).GetUserUnitList),
		beego.NSCtrlGet("/auth-menu/get-async-routes", (*adminAuth.MenuController).GetAsyncRoutes),

		beego.NSCtrlPost("/upload/upload", (*adminUpload.UploadController).Upload),
		beego.NSCtrlGet("/upload/vue-slice-upload", (*adminUpload.UploadController).VueSliceUploadCheck),
		beego.NSCtrlPost("/upload/vue-slice-upload", (*adminUpload.UploadController).VueSliceUpload),

		// system
		beego.NSCtrlGet("/system-unit/get", (*adminSystem.UnitController).GetUnitList),
	)

	// 请求前、后处理
	ns.Filter("before", func(ctx *context.Context) {
		(new(middleware.AccessMiddleware).RouterBefore())(ctx)
		middleware.AuthAdmin(&platWhiteApiList, &platAuthApiList)(ctx)
	})
	// ns.Filter("after", func(ctx *context.Context) {
	// 	(new(middleware.AccessMiddleware).RouterAfter())(ctx) // 请求后处理存在bug
	// })
	beego.InsertFilter("/admin_plat/*", beego.FinishRouter, new(middleware.AccessMiddleware).RouterAfter(&platWhiteApiList, &platAuthApiList), beego.WithReturnOnOutput(false))

	beego.AddNamespace(ns)
}
