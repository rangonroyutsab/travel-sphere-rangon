package api

import beego "github.com/beego/beego/v2/server/web"

type DashboardSummaryController struct {
	beego.Controller
}

func (c *DashboardSummaryController) Get() {
	c.Data["json"] = map[string]int{
		"total_saved": 0,
		"planned":     0,
		"visited":     0,
	}
	c.ServeJSON()
}
