package helper

import (
	"WenBeego/apps/common/dto"
	"encoding/json"
	"errors"
	"strings"

	"github.com/beego/beego/v2/server/web"
	beecontext "github.com/beego/beego/v2/server/web/context"
)

// 获取请求的token
func GetReqToken(ctx beecontext.Context) string {
	token := ctx.Request.Header.Get("Authorization")
	token = strings.TrimPrefix(token, "Bearer ")
	return token
}

// 获取请求接口
func GetReqUrl(ctx beecontext.Context) string {
	return ctx.Request.URL.Path
}

// 获取请求方法
func GetReqMethod(ctx beecontext.Context) string {
	return ctx.Request.Method
}

// 获取请求协议
func GetReqScheme(ctx beecontext.Context) string {
	return ctx.Request.URL.Scheme
}

// 获取请求主机
func GetReqHost(ctx beecontext.Context) string {
	return ctx.Request.Host
}

func GetReqIp(ctx beecontext.Context) string {
	return ctx.Request.RemoteAddr
}

func GetReqBody[T any](ctx *beecontext.Context) (T, error) {
	var data T
	reqBody := ctx.Input.RequestBody
	if len(reqBody) > 0 {
		err := json.Unmarshal(reqBody, &data)
		if err != nil {
			return data, err
		}
	}
	return data, nil
}

// 设置并获取基本参数
func GetBaseParamDto(ctx *beecontext.Context, moduleName string) (dto.BaseParamDto, error) {
	host := ctx.Request.Host
	tmpUserId := ctx.Input.GetData("userId")
	tmpUnitId := ctx.Input.GetData("unitId")
	tmpUnitUserId := ctx.Input.GetData("unitUserId")
	tmpIsOfficial := ctx.Input.GetData("isOfficial")
	userId := tmpUserId.(string)
	unitId := tmpUnitId.(string)
	unitUserId := tmpUnitUserId.(string)
	isOfficial := tmpIsOfficial.(bool)

	data := dto.BaseParamDto{}
	if moduleName == "" {
		return data, errors.New("moduleName is empty")
	}
	if unitId == "" {
		return data, errors.New("unitId is empty")
	}
	if userId == "" {
		return data, errors.New("userId is empty")
	}
	data.Host = host
	data.ModuleName = moduleName
	data.UnitId = unitId
	data.UserId = userId
	data.UnitUserId = unitUserId
	data.IsOfficial = isOfficial
	return data, nil
}

// 设置并返回请求页面数据列表参数
func GetReqDataListDto(ctxCtrl *web.Controller) (dto.ReqDataListDto, error) {
	data := dto.ReqDataListDto{}

	currentPage, err1 := ctxCtrl.GetInt("currentPage", 1)
	pageSize, err2 := ctxCtrl.GetInt("pageSize", 10)
	if err1 != nil {
		return data, err1
	}
	if err2 != nil {
		return data, err2
	}

	if pageSize <= 10 {
		pageSize = 10
	} else if pageSize > 10 && pageSize <= 20 {
		pageSize = 20
	} else {
		pageSize = 50
	}
	if currentPage <= 0 {
		currentPage = 1
	}

	data.PageSize = pageSize
	data.CurrentPage = currentPage
	data.Offset = (currentPage - 1) * pageSize
	return data, nil
}
