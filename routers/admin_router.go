package routers

import (
	adminAuth "WenBeego/apps/admin/controllers/auth"
	"WenBeego/apps/common/middleware"

	beego "github.com/beego/beego/v2/server/web"
	"github.com/beego/beego/v2/server/web/context"
)

func init() {
	ns := beego.NewNamespace("/admin",
		beego.NSCtrlGet("/login", (*adminAuth.IndexController).Get),
	)

	// 请求前、后处理
	ns.Filter("before", func(ctx *context.Context) {
		(new(middleware.AccessMiddleware).RouterBefore())(ctx)
	})
	// ns.Filter("after", func(ctx *context.Context) {
	// 	(new(middleware.AccessMiddleware).RouterAfter())(ctx) // 请求后处理存在bug
	// })
	beego.InsertFilter("/admin/*", beego.FinishRouter, new(middleware.AccessMiddleware).RouterAfter(), beego.WithReturnOnOutput(false))

	beego.AddNamespace(ns)
}
