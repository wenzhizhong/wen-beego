package auth

import (
	// beego "github.com/beego/beego/v2/server/web"

	"WenBeego/apps/admin_plat/services"
	commonControllers "WenBeego/apps/common/controller"
	"WenBeego/apps/common/helper"

	"WenBeego/apps/common/dto"
)

type AuthController struct {
	// beego.Controller
	commonControllers.AdminBaseController
	AuthService services.Auth
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
// @Router       /admin_plat/auth/login [post]
func (c *AuthController) Login() {
	loginDto, err := helper.GetReqBody[dto.LoginDto](c.Ctx)
	if err != nil {
		c.Data["json"] = helper.Response{Code: 0, Message: err.Error()}
		c.ServeJSON()
		return
	}

	data, err := c.AuthService.Login(loginDto, c.ModuleName)
	if err != nil {
		c.Data["json"] = helper.Response{Code: 0, Message: err.Error()}
		c.ServeJSON()
		return
	}

	c.Data["json"] = helper.Response{Code: 200, Message: "登录成功", Data: data}
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
// @Router       /admin_plat/auth/logout [post]
// @Security     ApiKeyAuth
func (c *AuthController) Logout() {

}

func (c *AuthController) Register() {
}

func (c *AuthController) ForgetPassword() {
}
func (c *AuthController) ResetPassword() {
}

func (c *AuthController) GetUserInfo() {
}
func (c *AuthController) UpdateUserInfo() {
}

// 获取验证码
// @Summary      获取验证码
// @Description  获取验证码
// @Tags         admin
// @Accept       json
// @Produce      json
// @Success      200  {object}  helper.Response
// @Failure      0  {object}  helper.Response
// @Router       /admin_plat/auth/get-captcha [get]
// @Security     ApiKeyAuth
func (c *AuthController) GetCatpcha() {
	data, err := c.AuthService.GetCatpcha()
	if err != nil {
		c.Data["json"] = helper.Response{Code: 0, Message: err.Error()}
		c.ServeJSON()
		return
	}

	c.Data["json"] = helper.Response{Code: 200, Message: "获取验证码成功", Data: data}
	c.ServeJSON()
}
