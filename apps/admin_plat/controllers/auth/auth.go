package auth

//身份认证-登录/注册/退出/获取用户信息等
import (
	// beego "github.com/beego/beego/v2/server/web"

	authService "WenBeego/apps/admin_plat/services/auth"
	commonControllers "WenBeego/apps/common/controller"
	"WenBeego/apps/common/dto/auth_dto"
	"WenBeego/apps/common/helper"
)

type AuthController struct {
	// beego.Controller
	commonControllers.AdminBaseController
	AuthService authService.Auth
}

// 登录
// @Summary      登录
// @Description  admin用户登录
// @Tags         admin
// @Accept       json
// @Produce      json
// @Param data body auth_dto.LoginDto true "登录参数"
// @Success      200  {object}  dto.Response
// @Failure      0  {object}  dto.Response
// @Router       /admin_plat/auth/login [post]
func (c *AuthController) Login() {
	host := c.Ctx.Request.Host
	loginDto, err := helper.GetReqBody[auth_dto.LoginDto](c.Ctx)
	if err != nil {
		c.Data["json"] = helper.Response(0, err.Error(), nil)
		c.ServeJSON()
		return
	}

	data, err := c.AuthService.Login(loginDto, c.ModuleName, host)
	if err != nil {
		c.Data["json"] = helper.Response(0, err.Error(), nil)
		c.ServeJSON()
		return
	}

	c.Data["json"] = helper.Response(200, "登录成功", data)
	c.ServeJSON()
}

// 退出登录
// @Summary      退出登录
// @Description  admin退出登录
// @Tags         admin
// @Accept       json
// @Produce      json
// @Success      200  {object}  dto.Response
// @Failure      0  {object}  dto.Response
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
// @Success      200  {object}  dto.Response
// @Failure      0  {object}  dto.Response
// @Router       /admin_plat/auth/get-captcha [get]
// @Security     ApiKeyAuth
func (c *AuthController) GetCatpcha() {
	data, err := c.AuthService.GetCatpcha()
	if err != nil {
		c.Data["json"] = helper.Response(0, err.Error(), nil)
		c.ServeJSON()
		return
	}

	c.Data["json"] = helper.Response(200, "获取验证码成功", data)
	c.ServeJSON()
}

// 刷新token
// @Summary      刷新token
// @Description  刷新token
// @Tags         admin
// @Accept       json
// @Produce      json
// @Success      200  {object}  dto.Response
// @Failure      0  {object}  dto.Response
// @Router       /admin_plat/auth/refresh-token [post]
// @Security     ApiKeyAuth
func (c *AuthController) RefreshToken() {
	body, err := helper.GetReqBody[auth_dto.RefreshTokenDto](c.Ctx)
	if err != nil {
		c.Data["json"] = helper.Response(0, err.Error(), nil)
		c.ServeJSON()
		return
	}

	data, err := c.AuthService.RefreshToken(c.ModuleName, body.BrancaToken, body.RefreshToken)
	if err != nil {
		c.Data["json"] = helper.Response(0, err.Error(), nil)
		c.ServeJSON()
		return
	}

	c.Data["json"] = helper.Response(200, "刷新成功", data)
	c.ServeJSON()
}
