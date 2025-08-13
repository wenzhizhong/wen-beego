package controller

type IndexBaseController struct {
	BaseController
}

func (c *IndexBaseController) Prepare() {
	c.BaseController.Prepare()
}
