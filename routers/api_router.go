package routers

import (
	apiAuthV1 "WenBeego/apps/api/controllers/auth/v1"
	apiUploadV1 "WenBeego/apps/api/controllers/upload/v1"
	"WenBeego/apps/common/middleware"

	beego "github.com/beego/beego/v2/server/web"
	"github.com/beego/beego/v2/server/web/context"
)

// 白名单api
var apiWhiteApiList = []string{
	"/api/v1/auth/login",
	"/api/v1/auth/register",
	"/api/v1/auth/forget-password",
	"/api/v1/auth/get-captcha",
	"/api/v1/auth/refresh-token",
}

// 登录后基础api
var apiAuthApiList = []string{
	"/api/v1/auth/logout",
	"/api/v1/upload/upload",
	"/api/v1/upload/vue-slice-upload",
	"/api/v1/upload/link-sign",
}

func apiAuthSlices() []beego.LinkNamespace {
	return []beego.LinkNamespace{
		// auth begin
		beego.NSCtrlPost("/v1/auth/login", (*apiAuthV1.AuthController).Login),
		beego.NSCtrlGet("/v1/auth/get-captcha", (*apiAuthV1.AuthController).GetCatpcha),
		beego.NSCtrlPost("/v1/auth/refresh-token", (*apiAuthV1.AuthController).RefreshToken),

		beego.NSCtrlPost("/v1/upload/upload", (*apiUploadV1.UploadController).Upload),
		beego.NSCtrlGet("/v1/upload/vue-slice-upload", (*apiUploadV1.UploadController).VueSliceUploadCheck),
		beego.NSCtrlPost("/v1/upload/vue-slice-upload", (*apiUploadV1.UploadController).VueSliceUpload),
		beego.NSCtrlGet("/v1/upload/link-sign", (*apiUploadV1.UploadController).LinkSign),
		// auth end
	}
}
func init() {
	allNamespaces := apiAuthSlices()
	ns := beego.NewNamespace("/api",
		allNamespaces...,
	)

	// 请求前、后处理
	ns.Filter("before", func(ctx *context.Context) {
		(new(middleware.AccessMiddleware).RouterBefore())(ctx)
		middleware.AuthUser(&apiWhiteApiList, &apiAuthApiList)(ctx)
	})
	// ns.Filter("after", func(ctx *context.Context) {
	// 	(new(middleware.AccessMiddleware).RouterAfter())(ctx) // 请求后处理存在bug
	// })
	beego.InsertFilter("/api/*", beego.FinishRouter, new(middleware.AccessMiddleware).RouterAfter(&apiWhiteApiList, &apiAuthApiList), beego.WithReturnOnOutput(false))

	beego.AddNamespace(ns)
}
