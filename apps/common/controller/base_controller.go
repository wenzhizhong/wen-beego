package controller

import (
	"strings"

	"github.com/beego/beego/v2/server/web"
	"github.com/beego/beego/v2/server/web/context" // 正确导入context包
)

type BaseController struct {
	web.Controller
	ModuleName string
}

func (c *BaseController) Prepare() {
	// 从路由中解析模块名
	c.ModuleName = parseModuleFromRoute(c.Ctx) // 使用正确的context类型
	if c.ModuleName == "" {
		c.ModuleName = "index"
	}
	// 设置视图路径
	c.ViewPath = "apps/" + c.ModuleName + "/views"
}

// parseModuleFromRoute 从路由路径解析模块名
func parseModuleFromRoute(ctx *context.Context) string {
	// 获取当前路由路径
	path := ctx.Request.URL.Path

	// 移除开头的斜杠
	if strings.HasPrefix(path, "/") {
		path = path[1:]
	}

	// 分割路径部分
	parts := strings.Split(path, "/")

	// 路径格式: /模块/控制器/方法
	if len(parts) >= 1 {
		return parts[0]
	}

	return ""
}
