package middleware

import (
	"WenBeego/apps/common/helper"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/beego/beego/v2/server/web"
	beecontext "github.com/beego/beego/v2/server/web/context"
)

func Auth(whiteApiList *[]string, authApiList *[]string) web.FilterFunc {
	return func(ctx *beecontext.Context) {
		tmpWhiteApiListMap := listToMap(*whiteApiList)
		tmpAuthApiListMap := listToMap(*authApiList)

		// 验证:白名单api
		isWriteApi, err := inWhiteApiList(ctx, tmpWhiteApiListMap)
		if err != nil {
			return
		}
		if isWriteApi {
			return
		}

		// 验证:认证token是否有效
		brancaData, err := checkToken(ctx)
		if err != nil {
			return
		}
		ctx.Input.SetData("userId", brancaData.Sub)

		// 验证:认证后基础api,通过则不校验权限
		isValid := checkAuth(ctx, brancaData.Sub, tmpAuthApiListMap)
		if isValid {
			return
		}

		// 检测用户是否有api权限
		checkUrlPermis(ctx, brancaData.Sub)
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
func inWhiteApiList(ctx *beecontext.Context, writeApiListMap map[string]bool) (bool, error) {
	url := helper.GetReqUrl(*ctx)
	if writeApiListMap[url] {
		return true, nil
	}

	jsonString := resposeStr(http.StatusUnauthorized, "请登录！", nil)
	ctx.ResponseWriter.ResponseWriter.WriteHeader(http.StatusUnauthorized)
	ctx.ResponseWriter.Write([]byte(jsonString))
	return false, errors.New("请登录！")
}

// 验证:登录后基本api
func inAuthApiList(ctx *beecontext.Context, authApiListMap map[string]bool) bool {
	url := helper.GetReqUrl(*ctx)
	if authApiListMap[url] {
		return true
	}
	return false
}

// 校验token
func checkToken(ctx *beecontext.Context) (helper.BrancaData, error) {
	token := helper.GetReqToken(*ctx)
	brancaData, err := helper.BrancaDecode(token)
	if err != nil {
		jsonString := resposeStr(http.StatusUnauthorized, err.Error(), nil)
		ctx.ResponseWriter.ResponseWriter.WriteHeader(http.StatusUnauthorized)
		ctx.ResponseWriter.Write([]byte(jsonString))
		return brancaData, err
	}
	return brancaData, nil
}

// 校验登录后基本路由，通过则不校验权限
func checkAuth(ctx *beecontext.Context, userId string, authApiListMap map[string]bool) bool {
	hasBaseAuthApi := inAuthApiList(ctx, authApiListMap)
	if hasBaseAuthApi {
		return true
	}
	return false
}

// 检测用户是否有api权限
func checkUrlPermis(ctx *beecontext.Context, userId string) bool {
	hasPermis := helper.HasUrlPermis(userId, helper.GetReqUrl(*ctx))
	if !hasPermis {
		jsonString := resposeStr(http.StatusUnauthorized, "没有接口权限", nil)
		ctx.ResponseWriter.ResponseWriter.WriteHeader(http.StatusUnauthorized)
		ctx.ResponseWriter.Write([]byte(jsonString))
		return false
	}
	return hasPermis
}
