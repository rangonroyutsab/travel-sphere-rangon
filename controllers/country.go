package controllers

type CountryController struct {
	BaseController
}

func (c *CountryController) Get() {
	c.Data["Title"] = "Countries"
	c.TplName = "countries.tpl"
}

type CountryDetailController struct {
	BaseController
}

func (c *CountryDetailController) Get() {
	c.Data["Title"] = "Destination"
	c.TplName = "destination.tpl"
}
