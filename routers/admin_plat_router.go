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
	"/admin_plat/system-dept/get-dept-tree",
	"/admin_plat/system-dept/get-dept-principal",
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
		beego.NSCtrlGet("/system-unit/get", (*adminSystem.UnitController).Get),
		beego.NSCtrlPost("/system-unit/add", (*adminSystem.UnitController).Add),
		beego.NSCtrlPost("/system-unit/edit", (*adminSystem.UnitController).Edit),
		beego.NSCtrlPost("/system-unit/del", (*adminSystem.UnitController).Del),

		beego.NSCtrlGet("/system-dept/get", (*adminSystem.DeptController).Get),
		beego.NSCtrlPost("/system-dept/add", (*adminSystem.DeptController).Add),
		beego.NSCtrlPost("/system-dept/edit", (*adminSystem.DeptController).Edit),
		beego.NSCtrlPost("/system-dept/del", (*adminSystem.DeptController).Del),
		beego.NSCtrlGet("/system-dept/get-dept-tree", (*adminSystem.DeptController).GetUnitDeptTree),
		beego.NSCtrlGet("/system-dept/get-dept-principal", (*adminSystem.DeptController).GetUnitDeptPrincipal),

		beego.NSCtrlGet("/system-role/get", (*adminSystem.RoleController).Get),
		beego.NSCtrlPost("/system-role/add", (*adminSystem.RoleController).Add),
		beego.NSCtrlPost("/system-role/edit", (*adminSystem.RoleController).Edit),
		beego.NSCtrlPost("/system-role/del", (*adminSystem.RoleController).Del),
		beego.NSCtrlGet("/system-role/role-menu", (*adminSystem.RoleController).RoleMenu),
		beego.NSCtrlGet("/system-role/role-menu-ids", (*adminSystem.RoleController).RoleMenuIds),
		beego.NSCtrlPost("/system-role/role-menu-save", (*adminSystem.RoleController).RoleMenuSave),

		beego.NSCtrlGet("/system-user/get", (*adminSystem.UserController).GetUserList),
		beego.NSCtrlPost("/system-user/add", (*adminSystem.UserController).Add),
		beego.NSCtrlPost("/system-user/edit", (*adminSystem.UserController).Edit),
		beego.NSCtrlPost("/system-user/del", (*adminSystem.UserController).Del),
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
