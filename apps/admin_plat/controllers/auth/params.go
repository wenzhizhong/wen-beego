package auth

import (
	commonControllers "WenBeego/apps/common/controller"
	"WenBeego/apps/common/helper"
	commonAuthService "WenBeego/apps/common/services/auth"
)

type ParamsController struct {
	commonControllers.AdminBaseController
}

// 获取参数-模型常量参数
// @Summary 获取参数-模型常量参数
// @Description 获取参数-模型常量参数
// @Tags 参数
// @Success 200 {object} dto.Response "返回结果"
// @Router /admin_plat/auth-params/model-params [get]
func (c *ParamsController) GetModelParams() {
	data := map[string]interface{}{
		"modelParam": commonAuthService.GetAllModelConstant(),
	}
	c.Data["json"] = helper.Response(200, "刷新成功", data)
	c.ServeJSON()
}
