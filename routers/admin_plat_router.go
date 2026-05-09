package routers

import (
	adminAuth "WenBeego/apps/admin_plat/controllers/auth"
	adminMonitor "WenBeego/apps/admin_plat/controllers/monitor"
	adminSystem "WenBeego/apps/admin_plat/controllers/system"
	adminUpload "WenBeego/apps/admin_plat/controllers/upload"
	"WenBeego/apps/common/global/constant"
	"WenBeego/apps/common/middleware/blocker"

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
	"/admin_plat/upload/link-sign",
}

func platAuthSlices() []beego.LinkNamespace {
	return []beego.LinkNamespace{
		// auth begin
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
		beego.NSCtrlGet("/upload/link-sign", (*adminUpload.UploadController).LinkSign),
		// auth end
	}
}
func platSystemSlices() []beego.LinkNamespace {
	return []beego.LinkNamespace{
		// system begin
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

		beego.NSCtrlGet("/system-menu/get-plat", (*adminSystem.MenuPlatController).Get),
		beego.NSCtrlPost("/system-menu/add-plat", (*adminSystem.MenuPlatController).Add),
		beego.NSCtrlPost("/system-menu/edit-plat", (*adminSystem.MenuPlatController).Edit),
		beego.NSCtrlPost("/system-menu/del-plat", (*adminSystem.MenuPlatController).Del),
		beego.NSCtrlGet("/system-menu/get-mchnt", (*adminSystem.MenuMchntController).Get),
		beego.NSCtrlPost("/system-menu/add-mchnt", (*adminSystem.MenuMchntController).Add),
		beego.NSCtrlPost("/system-menu/edit-mchnt", (*adminSystem.MenuMchntController).Edit),
		beego.NSCtrlPost("/system-menu/del-mchnt", (*adminSystem.MenuMchntController).Del),
		beego.NSCtrlGet("/system-menu/mchnt-unit-list", (*adminSystem.MenuMchntController).MchntUnitList),

		beego.NSCtrlGet("/monitor-cron/get", (*adminMonitor.CronController).Get),
		beego.NSCtrlPost("/monitor-cron/add", (*adminMonitor.CronController).Add),
		beego.NSCtrlPost("/monitor-cron/edit", (*adminMonitor.CronController).Edit),
		beego.NSCtrlPost("/monitor-cron/del", (*adminMonitor.CronController).Del),
		beego.NSCtrlPost("/monitor-cron/change-status", (*adminMonitor.CronController).ChangeStatus),
		beego.NSCtrlGet("/monitor-cron/get-avaible", (*adminMonitor.CronController).GetAvaibleCronList),
		beego.NSCtrlGet("/monitor-cron-log/get", (*adminMonitor.CronLogController).Get),

		beego.NSCtrlGet("/monitor-queue-dlq/get", (*adminMonitor.QueueDlqController).Get),
		beego.NSCtrlPost("/monitor-queue-dlq/requeue", (*adminMonitor.QueueDlqController).Requeue),
		// system end
	}
}

func init() {
	allNamespaces := append(commonSlices(constant.ADMIN_PLAT), platAuthSlices()...)
	allNamespaces = append(allNamespaces, platSystemSlices()...)
	ns := beego.NewNamespace("/admin_plat",
		allNamespaces...,
	)

	// 请求前、后处理
	ns.Filter("before", func(ctx *context.Context) {
		blocker.RouterBefore()(ctx)
		blocker.AuthAdmin(&platWhiteApiList, &platAuthApiList)(ctx)
	})
	// ns.Filter("after", func(ctx *context.Context) {
	// 	(blocker.RouterAfter())(ctx) // 请求后处理存在bug
	// })
	beego.InsertFilter("/admin_plat/*", beego.FinishRouter, blocker.RouterAfter(&platWhiteApiList, &platAuthApiList), beego.WithReturnOnOutput(false))

	beego.AddNamespace(ns)
}
