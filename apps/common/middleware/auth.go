package middleware

import (
	"WenBeego/apps/common/global"
	"WenBeego/apps/common/helper"
	"WenBeego/apps/common/services/framework"
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
		moduleName := helper.ParseModuleFromRoute(ctx.Request.URL.Path)
		brancaData, err := checkToken(ctx, moduleName)
		if err != nil {
			return
		}
		// 验证:是否是官方平台
		isOfficial := false
		isOfficial, err = helper.IsOfficial(moduleName, brancaData.SubUnit)
		if err != nil {
			return
		}
		ctx.Input.SetData("userId", brancaData.Sub)
		ctx.Input.SetData("unitUserId", brancaData.SubUnitUser) // （plat_user表/mchnt_user表，用户所在单位组织）
		ctx.Input.SetData("unitId", brancaData.SubUnit)
		ctx.Input.SetData("isOfficial", isOfficial)

		// 验证:认证后基础api,通过则不校验权限
		isValid := checkBaseAuthApi(ctx, tmpAuthApiListMap)
		if isValid {
			return
		}

		type checkResult struct {
			status    bool
			err       error
			checkType string
		}
		checkResultNumber := 2
		ch := make(chan checkResult, checkResultNumber)
		// 验证状态:用户状态/组织单位状态/角色状态
		go func() {
			status, err := checkAuthAdminStatus(moduleName, brancaData)
			ch <- checkResult{status: status, err: err, checkType: "status"}
		}()
		// 检测用户是否有api权限
		go func() {
			status, err := checkAuthAdminPermis(moduleName, brancaData, helper.GetReqUrl(*ctx))
			ch <- checkResult{status: status, err: err, checkType: "perms"}
		}()
		errorStr := ""
		for i := 0; i < checkResultNumber; i++ {
			result := <-ch
			if result.err != nil {
				errorStr += result.err.Error() + ";\n"
			}
		}
		close(ch)
		if errorStr != "" {
			setResponse(ctx, http.StatusNetworkAuthenticationRequired, errorStr, nil)
			return
		}
	}
}
func responseStr(code int, msg string, data interface{}) string {
	res := helper.Response(code, msg, data)
	jsonString, _ := json.Marshal(res)
	return string(jsonString)
}
func setResponse(ctx *beecontext.Context, code int, msg string, data interface{}) {
	jsonString := responseStr(code, msg, data)
	ctx.ResponseWriter.ResponseWriter.WriteHeader(code)
	ctx.ResponseWriter.Write([]byte(jsonString))
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
		setResponse(ctx, http.StatusUnauthorized, err.Error(), nil)
		return brancaData, err
	}

	aud, _ := global.GetConfigDiy("branca." + moduleName + ".aud")
	iss, _ := global.GetConfigDiy("branca." + moduleName + ".iss")
	if aud != nil && iss != nil && (brancaData.Aud != aud.(string) || brancaData.Iss != iss.(string)) {
		setResponse(ctx, http.StatusUnauthorized, "无效的token", nil)
		return brancaData, err
	}
	return brancaData, nil
}

// 校验登录后基本路由，通过则不校验权限
func checkBaseAuthApi(ctx *beecontext.Context, authApiListMap map[string]bool) bool {
	hasBaseAuthApi := inAuthApiList(ctx, authApiListMap)
	return hasBaseAuthApi
}

// 检测用户是否有api权限
func checkAuthAdminPermis(moduleName string, brancaData helper.BrancaData, path string) (bool, error) {
	service := &framework.AuthMiddlewate{}
	status, err := service.CheckAuthAdminRouters(moduleName, brancaData, path)
	return status, err
}

// 检测用户状态
func checkAuthAdminStatus(moduleName string, brancaData helper.BrancaData) (bool, error) {
	service := &framework.AuthMiddlewate{}
	status, err := service.CheckAuthAdminStatus(moduleName, brancaData)
	return status, err
}
