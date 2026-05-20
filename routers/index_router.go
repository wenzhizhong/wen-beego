package routers

import (
	"WenBeego/apps/common/middleware/blocker"
	indeHome "WenBeego/apps/index/controllers/home"

	beego "github.com/beego/beego/v2/server/web"
	"github.com/beego/beego/v2/server/web/context"
)

var indexWhiteApiList, indexAuthApiList []string

func init() {
	// beego.Router("/", &indeHome.IndexController{})
	// beego.Router("/index", &indeHome.IndexController{})
	// beego.Router("/index/index", &indeHome.IndexController{})

	beego.Get("/", func(ctx *context.Context) {
		ctx.Redirect(302, "/index")
	})
	ns := beego.NewNamespace("/index",
		beego.NSCtrlGet("/", (*indeHome.IndexController).Get),
		beego.NSCtrlGet("/index", (*indeHome.IndexController).Get),
		beego.NSCtrlPost("/index", (*indeHome.IndexController).Post),
	)

	// 请求前、后处理
	ns.Filter("before", func(ctx *context.Context) {
		blocker.RouterBefore()(ctx)
		if ctx.ResponseWriter.Started {
			return
		}
	})
	// ns.Filter("after", func(ctx *context.Context) {
	// 	(blocker.RouterAfter())(ctx) // 请求后处理存在bug
	// })
	beego.InsertFilter("/index/*", beego.FinishRouter, blocker.RouterAfter(&indexWhiteApiList, &indexAuthApiList), beego.WithReturnOnOutput(false))

	beego.AddNamespace(ns)
}
