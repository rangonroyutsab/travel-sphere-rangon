package controllers

import (
	"strings"

	beego "github.com/beego/beego/v2/server/web"
)

type LoginController struct {
	BaseController
}

func (c *LoginController) Get() {
	c.Data["Title"] = "Login"
	c.TplName = "login.tpl"
}

func (c *LoginController) Post() {
	username := strings.TrimSpace(c.GetString("username"))
	password := c.GetString("password")

	demoUsername, _ := beego.AppConfig.String("DEMO_USERNAME")
	demoPassword, _ := beego.AppConfig.String("DEMO_PASSWORD")

	if username == demoUsername && password == demoPassword {
		c.SetSession("username", username)
		c.Redirect("/", 302)
		return
	}

	c.Data["Title"] = "Login"
	c.Data["Error"] = "Invalid username or password"
	c.TplName = "login.tpl"
}

type LogoutController struct {
	BaseController
}

func (c *LogoutController) Get() {
	c.DelSession("username")
	c.Redirect("/", 302)
}
