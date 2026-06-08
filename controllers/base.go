package controllers

import (
	beego "github.com/beego/beego/v2/server/web"
)

type BaseController struct {
	beego.Controller
}

func (c *BaseController) Prepare() {
	c.Layout = "layout.tpl"

	c.Data["AppName"] = "TravelSphere"
	c.Data["Title"] = "TravelSphere"
	c.Data["CurrentPath"] = c.Ctx.Request.URL.Path

	username, ok := c.GetSession("username").(string)
	if ok && username != "" {
		c.Data["IsAuthenticated"] = true
		c.Data["UserName"] = username
		return
	}

	c.Data["IsAuthenticated"] = false
	c.Data["UserName"] = ""
}

func (c *BaseController) CurrentUsername() string {
	username, _ := c.GetSession("username").(string)
	return username
}
