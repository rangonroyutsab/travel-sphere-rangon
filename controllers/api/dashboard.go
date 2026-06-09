package api

import (
	"travel-sphere-rangon/services"

	beego "github.com/beego/beego/v2/server/web"
)

type DashboardSummaryController struct {
	beego.Controller
}

func (c *DashboardSummaryController) Get() {
	username := currentUser(&c.Controller)
	summary := services.Dashboard.GetSummary(username)

	c.Data["json"] = summary
	c.ServeJSON()
}
