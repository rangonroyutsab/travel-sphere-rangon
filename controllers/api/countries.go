package api

import beego "github.com/beego/beego/v2/server/web"

type CountriesController struct {
	beego.Controller
}

func (c *CountriesController) Get() {
	c.Data["json"] = []string{}
	c.ServeJSON()
}

type CountryDetailController struct {
	beego.Controller
}

func (c *CountryDetailController) Get() {
	c.Data["json"] = map[string]string{
		"message": "country detail placeholder",
	}
	c.ServeJSON()
}
