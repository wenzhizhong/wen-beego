package routers

import (
	adminAuth "WenBeego/apps/admin_mchnt/controllers/auth"
	adminSystem "WenBeego/apps/admin_mchnt/controllers/system"
	adminUpload "WenBeego/apps/admin_mchnt/controllers/upload"
	adminSystemPlat "WenBeego/apps/admin_plat/controllers/system_mchnt"
	"WenBeego/apps/common/global/constant"
	"WenBeego/apps/common/middleware/blocker"

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

func commonSlices(gotoModuleName string) []beego.LinkNamespace {
	unit := []beego.LinkNamespace{
		// system begin
		beego.NSCtrlGet("/system-unit/get", (*adminSystem.UnitController).Get),
		beego.NSCtrlPost("/system-unit/add", (*adminSystem.UnitController).Add),
		beego.NSCtrlPost("/system-unit/edit", (*adminSystem.UnitController).Edit),
		beego.NSCtrlPost("/system-unit/del", (*adminSystem.UnitController).Del),
	}
	if gotoModuleName == constant.ADMIN_PLAT {
		unit = []beego.LinkNamespace{
			beego.NSCtrlGet("/admin_mchnt/system-unit/get", (*adminSystemPlat.UnitController).Get),
			beego.NSCtrlPost("/admin_mchnt/system-unit/add", (*adminSystemPlat.UnitController).Add),
			beego.NSCtrlPost("/admin_mchnt/system-unit/edit", (*adminSystemPlat.UnitController).Edit),
			beego.NSCtrlPost("/admin_mchnt/system-unit/del", (*adminSystemPlat.UnitController).Del),
		}
	}

	dept := []beego.LinkNamespace{
		beego.NSCtrlGet("/system-dept/get", (*adminSystem.DeptController).Get),
		beego.NSCtrlPost("/system-dept/add", (*adminSystem.DeptController).Add),
		beego.NSCtrlPost("/system-dept/edit", (*adminSystem.DeptController).Edit),
		beego.NSCtrlPost("/system-dept/del", (*adminSystem.DeptController).Del),
		beego.NSCtrlGet("/system-dept/get-dept-tree", (*adminSystem.DeptController).GetUnitDeptTree),
		beego.NSCtrlGet("/system-dept/get-dept-principal", (*adminSystem.DeptController).GetUnitDeptPrincipal),
	}
	if gotoModuleName == constant.ADMIN_PLAT {
		dept = []beego.LinkNamespace{
			beego.NSCtrlGet("/admin_mchnt/system-dept/get", (*adminSystemPlat.DeptController).Get),
			beego.NSCtrlPost("/admin_mchnt/system-dept/add", (*adminSystemPlat.DeptController).Add),
			beego.NSCtrlPost("/admin_mchnt/system-dept/edit", (*adminSystemPlat.DeptController).Edit),
			beego.NSCtrlPost("/admin_mchnt/system-dept/del", (*adminSystemPlat.DeptController).Del),
			beego.NSCtrlGet("/admin_mchnt/system-dept/get-dept-tree", (*adminSystemPlat.DeptController).GetUnitDeptTree),
			beego.NSCtrlGet("/admin_mchnt/system-dept/get-dept-principal", (*adminSystemPlat.DeptController).GetUnitDeptPrincipal),
		}
	}

	role := []beego.LinkNamespace{
		beego.NSCtrlGet("/system-role/get", (*adminSystem.RoleController).Get),
		beego.NSCtrlPost("/system-role/add", (*adminSystem.RoleController).Add),
		beego.NSCtrlPost("/system-role/edit", (*adminSystem.RoleController).Edit),
		beego.NSCtrlPost("/system-role/del", (*adminSystem.RoleController).Del),
		beego.NSCtrlGet("/system-role/role-menu", (*adminSystem.RoleController).RoleMenu),
		beego.NSCtrlGet("/system-role/role-menu-ids", (*adminSystem.RoleController).RoleMenuIds),
		beego.NSCtrlPost("/system-role/role-menu-save", (*adminSystem.RoleController).RoleMenuSave),
	}
	if gotoModuleName == constant.ADMIN_PLAT {
		role = []beego.LinkNamespace{
			beego.NSCtrlGet("/admin_mchnt/system-role/get", (*adminSystemPlat.RoleController).Get),
			beego.NSCtrlPost("/admin_mchnt/system-role/add", (*adminSystemPlat.RoleController).Add),
			beego.NSCtrlPost("/admin_mchnt/system-role/edit", (*adminSystemPlat.RoleController).Edit),
			beego.NSCtrlPost("/admin_mchnt/system-role/del", (*adminSystemPlat.RoleController).Del),
			beego.NSCtrlGet("/admin_mchnt/system-role/role-menu", (*adminSystemPlat.RoleController).RoleMenu),
			beego.NSCtrlGet("/admin_mchnt/system-role/role-menu-ids", (*adminSystemPlat.RoleController).RoleMenuIds),
			beego.NSCtrlPost("/admin_mchnt/system-role/role-menu-save", (*adminSystemPlat.RoleController).RoleMenuSave),
		}
	}

	user := []beego.LinkNamespace{
		beego.NSCtrlGet("/system-user/get", (*adminSystem.UserController).GetUserList),
		beego.NSCtrlPost("/system-user/add", (*adminSystem.UserController).Add),
		beego.NSCtrlPost("/system-user/edit", (*adminSystem.UserController).Edit),
		beego.NSCtrlPost("/system-user/del", (*adminSystem.UserController).Del),
		beego.NSCtrlGet("/system-user/get-role-tree", (*adminSystem.UserController).GetRoleTree),
	}
	if gotoModuleName == constant.ADMIN_PLAT {
		user = []beego.LinkNamespace{
			beego.NSCtrlGet("/admin_mchnt/system-user/get", (*adminSystemPlat.UserController).GetUserList),
			beego.NSCtrlPost("/admin_mchnt/system-user/add", (*adminSystemPlat.UserController).Add),
			beego.NSCtrlPost("/admin_mchnt/system-user/edit", (*adminSystemPlat.UserController).Edit),
			beego.NSCtrlPost("/admin_mchnt/system-user/del", (*adminSystemPlat.UserController).Del),
			beego.NSCtrlGet("/admin_mchnt/system-user/get-role-tree", (*adminSystemPlat.UserController).GetRoleTree),
		}
	}

	commonSlice := []beego.LinkNamespace{}
	commonSlice = append(commonSlice, unit...)
	commonSlice = append(commonSlice, dept...)
	commonSlice = append(commonSlice, role...)
	commonSlice = append(commonSlice, user...)
	return commonSlice
}
func mchntAuthSlices() []beego.LinkNamespace {
	return []beego.LinkNamespace{
		// auth begin

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
		// auth end
	}
}

func init() {
	allNamespaces := append(commonSlices(""), mchntAuthSlices()...)
	ns := beego.NewNamespace("/admin_mchnt",
		allNamespaces...,
	)

	// 请求前、后处理
	ns.Filter("before", func(ctx *context.Context) {
		(new(blocker.AccessMiddleware).RouterBefore())(ctx)
		blocker.AuthAdmin(&mchntWhiteApiList, &mchntAuthApiList)(ctx)
	})
	// ns.Filter("after", func(ctx *context.Context) {
	// 	(new(blocker.AccessMiddleware).RouterAfter())(ctx) // 请求后处理存在bug
	// })
	beego.InsertFilter("/admin_mchnt/*", beego.FinishRouter, new(blocker.AccessMiddleware).RouterAfter(&mchntWhiteApiList, &mchntAuthApiList), beego.WithReturnOnOutput(false))

	beego.AddNamespace(ns)
}
