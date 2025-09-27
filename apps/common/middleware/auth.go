package middleware

import (
	"WenBeego/apps/common/global"
	"WenBeego/apps/common/helper"
	"WenBeego/apps/common/services"
	"encoding/json"
	"net/http"

	"github.com/beego/beego/v2/server/web"
	beecontext "github.com/beego/beego/v2/server/web/context"
)

func AuthAdmin(whiteApiList *[]string, authApiList *[]string) web.FilterFunc {
	return func(ctx *beecontext.Context) {
		tmpWhiteApiListMap := listToMap(*whiteApiList)
		tmpAuthApiListMap := listToMap(*authApiList)

		// 验证:白名单api
		isWriteApi := inWhiteApiList(ctx, tmpWhiteApiListMap)
		if isWriteApi {
			return
		}

		// 验证:认证token是否有效
		moduleName := helper.ParseModuleFromRoute(ctx)
		brancaData, err := checkToken(ctx, moduleName)
		if err != nil {
			return
		}
		ctx.Input.SetData("userId", brancaData.Sub)
		ctx.Input.SetData("unitId", brancaData.SubUnit)

		// 验证:认证后基础api,通过则不校验权限
		isValid := checkAuth(ctx, tmpAuthApiListMap)
		if isValid {
			return
		}
		// 验证状态:用户状态/组织单位状态/角色状态
		_, err = checkAuthAdminStatus(ctx, moduleName, brancaData)
		if err != nil {
			return
		}

		// 检测用户是否有api权限
		checkUrlPermis(ctx, moduleName, brancaData)
	}
}
func resposeStr(code int, msg string, data interface{}) string {
	res := helper.Response{
		Code:    code,
		Message: msg,
		Data:    data,
	}
	jsonString, _ := json.Marshal(res)
	return string(jsonString)
}

// list 转 map
func listToMap(listData []string) map[string]bool {
	tmpMap := make(map[string]bool)
	for _, v := range listData {
		tmpMap[v] = true
	}
	return tmpMap
}

// 验证:白名单api
func inWhiteApiList(ctx *beecontext.Context, writeApiListMap map[string]bool) bool {
	url := helper.GetReqUrl(*ctx)
	return writeApiListMap[url]
}

// 验证:登录后基本api
func inAuthApiList(ctx *beecontext.Context, authApiListMap map[string]bool) bool {
	url := helper.GetReqUrl(*ctx)
	return authApiListMap[url]
}

// 校验token
func checkToken(ctx *beecontext.Context, moduleName string) (helper.BrancaData, error) {
	token := helper.GetReqToken(*ctx)
	brancaData, err := helper.BrancaDecode(token, moduleName)
	if err != nil {
		jsonString := resposeStr(http.StatusUnauthorized, err.Error(), nil)
		ctx.ResponseWriter.ResponseWriter.WriteHeader(http.StatusUnauthorized)
		ctx.ResponseWriter.Write([]byte(jsonString))
		return brancaData, err
	}

	aud, _ := global.GetConfigDiy("branca." + moduleName + ".aud")
	iss, _ := global.GetConfigDiy("branca." + moduleName + ".iss")
	if aud != nil && iss != nil && (brancaData.Aud != aud.(string) || brancaData.Iss != iss.(string)) {
		jsonString := resposeStr(http.StatusUnauthorized, "无效的token", nil)
		ctx.ResponseWriter.ResponseWriter.WriteHeader(http.StatusUnauthorized)
		ctx.ResponseWriter.Write([]byte(jsonString))
		return brancaData, err
	}
	return brancaData, nil
}

// 校验登录后基本路由，通过则不校验权限
func checkAuth(ctx *beecontext.Context, authApiListMap map[string]bool) bool {
	hasBaseAuthApi := inAuthApiList(ctx, authApiListMap)
	return hasBaseAuthApi
}

// 检测用户是否有api权限
func checkUrlPermis(ctx *beecontext.Context, moduleName string, brancaData helper.BrancaData) (bool, error) {
	service := &services.AuthMiddlewate{}
	status, err := service.CheckAuthAdminRouters(moduleName, brancaData, helper.GetReqUrl(*ctx))
	if err != nil {
		jsonString := resposeStr(http.StatusBadGateway, err.Error(), nil)
		ctx.ResponseWriter.ResponseWriter.WriteHeader(http.StatusBadGateway)
		ctx.ResponseWriter.Write([]byte(jsonString))
	}
	return status, err
}

// 检测用户状态
func checkAuthAdminStatus(ctx *beecontext.Context, moduleName string, brancaData helper.BrancaData) (bool, error) {
	service := &services.AuthMiddlewate{}
	status, err := service.CheckAuthAdminStatus(moduleName, brancaData)
	if err != nil {
		jsonString := resposeStr(http.StatusBadGateway, err.Error(), nil)
		ctx.ResponseWriter.ResponseWriter.WriteHeader(http.StatusBadGateway)
		ctx.ResponseWriter.Write([]byte(jsonString))
	}
	return status, err
}
