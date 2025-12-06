package routers

import (
	adminAuth "WenBeego/apps/admin_mchnt/controllers/auth"
	adminSystem "WenBeego/apps/admin_mchnt/controllers/system"
	adminUpload "WenBeego/apps/admin_mchnt/controllers/upload"
	"WenBeego/apps/common/middleware"

	beego "github.com/beego/beego/v2/server/web"
	"github.com/beego/beego/v2/server/web/context"
)

// 白名单api
var mchntWhiteApiList = []string{
	"/admin_mchnt/auth/login",
	"/admin_mchnt/auth/register",
	"/admin_mchnt/auth/forget-password",
	"/admin_mchnt/auth/get-captcha",
	"/admin_mchnt/auth/refresh-token",
}

// 登录后基础api
var mchntAuthApiList = []string{
	"/admin_mchnt/auth/logout",
	"/admin_mchnt/auth-params/model-params",
	"/admin_mchnt/auth-mchnt/change-unit",
	"/admin_mchnt/auth-mchnt/get-user-unit-list",
	"/admin_mchnt/auth-menu/get-async-routes",
	"/admin_mchnt/upload/upload",
	"/admin_mchnt/upload/vue-slice-upload",
	"/admin_mchnt/upload/link-sign",
}

func init() {
	ns := beego.NewNamespace("/admin_mchnt",
		// auth
		beego.NSCtrlPost("/auth/login", (*adminAuth.AuthController).Login),
		beego.NSCtrlGet("/auth/get-captcha", (*adminAuth.AuthController).GetCatpcha),
		beego.NSCtrlPost("/auth/refresh-token", (*adminAuth.AuthController).RefreshToken),
		beego.NSCtrlGet("/auth-params/model-params", (*adminAuth.ParamsController).GetModelParams),
		beego.NSCtrlPost("/auth-mchnt/change-unit", (*adminAuth.MchntController).ChangeUnit),
		beego.NSCtrlGet("/auth-mchnt/get-user-unit-list", (*adminAuth.MchntController).GetUserUnitList),
		beego.NSCtrlGet("/auth-menu/get-async-routes", (*adminAuth.MenuController).GetAsyncRoutes),

		beego.NSCtrlPost("/upload/upload", (*adminUpload.UploadController).Upload),
		beego.NSCtrlGet("/upload/vue-slice-upload", (*adminUpload.UploadController).VueSliceUploadCheck),
		beego.NSCtrlPost("/upload/vue-slice-upload", (*adminUpload.UploadController).VueSliceUpload),
		beego.NSCtrlGet("/upload/link-sign", (*adminUpload.UploadController).LinkSign),

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
		beego.NSCtrlGet("/system-user/get-role-tree", (*adminSystem.UserController).GetRoleTree),
	)

	// 请求前、后处理
	ns.Filter("before", func(ctx *context.Context) {
		(new(middleware.AccessMiddleware).RouterBefore())(ctx)
		middleware.AuthAdmin(&mchntWhiteApiList, &mchntAuthApiList)(ctx)
	})
	// ns.Filter("after", func(ctx *context.Context) {
	// 	(new(middleware.AccessMiddleware).RouterAfter())(ctx) // 请求后处理存在bug
	// })
	beego.InsertFilter("/admin_mchnt/*", beego.FinishRouter, new(middleware.AccessMiddleware).RouterAfter(&mchntWhiteApiList, &mchntAuthApiList), beego.WithReturnOnOutput(false))

	beego.AddNamespace(ns)
}
