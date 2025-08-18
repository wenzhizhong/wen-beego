package auth

import (
	// beego "github.com/beego/beego/v2/server/web"

	commonControllers "WenBeego/apps/common/controller"
	"WenBeego/apps/common/helper"

	"WenBeego/apps/common/dto"
)

type AuthController struct {
	// beego.Controller
	commonControllers.AdminBaseController
}

// 登录
// @Summary      登录
// @Description  admin用户登录
// @Tags         admin
// @Accept       json
// @Produce      json
// @Param data body dto.LoginDto true "登录参数"
// @Success      200  {object}  helper.Response
// @Failure      0  {object}  helper.Response
// @Router       /admin/auth/login [post]
func (c *AuthController) Login() {
	loginDto, err := helper.GetReqBody[dto.LoginDto](c.Ctx)
	if err != nil {
		c.Data["json"] = helper.Response{Code: 0, Message: err.Error()}
		c.ServeJSON()
		return
	}
	// TODO

	c.Data["json"] = helper.Response{Code: 200, Message: "登录成功"}
	c.ServeJSON()
}

// 退出登录
// @Summary      退出登录
// @Description  admin退出登录
// @Tags         admin
// @Accept       json
// @Produce      json
// @Success      200  {object}  helper.Response
// @Failure      0  {object}  helper.Response
// @Router       /admin/auth/logout [post]
// @Security     ApiKeyAuth
func (c *AuthController) Logout() {

}
