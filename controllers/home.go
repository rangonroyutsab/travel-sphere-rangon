package controllers

type HomeController struct {
	BaseController
}

func (c *HomeController) Get() {
	c.Data["Title"] = "Home"
	c.TplName = "home.tpl"
}
