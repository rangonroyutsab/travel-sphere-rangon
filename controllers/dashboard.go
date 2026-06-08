package controllers

type DashboardController struct {
	BaseController
}

func (c *DashboardController) Get() {
	c.Data["Title"] = "Dashboard"
	c.TplName = "dashboard.tpl"
}
