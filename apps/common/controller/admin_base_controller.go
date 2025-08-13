package controller

type AdminBaseController struct {
	BaseController
}

func (c *AdminBaseController) Prepare() {
	c.BaseController.Prepare()
}
