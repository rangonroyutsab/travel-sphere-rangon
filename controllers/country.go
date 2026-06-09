package controllers

import "travel-sphere-rangon/services"

type CountryController struct {
	BaseController
}

func (c *CountryController) Get() {
	c.Data["Title"] = "Countries"

	countries, err := services.Countries.GetDefaultCountries()
	if err != nil {
		c.Data["Error"] = "Countries are temporarily unavailable."
		countries = nil
	}

	c.Data["Countries"] = countries
	c.TplName = "countries.tpl"
}

type CountryDetailController struct {
	BaseController
}

func (c *CountryDetailController) Get() {
	slug := c.Ctx.Input.Param(":slug")

	country, err := services.Countries.GetCountryBySlug(slug)
	if err != nil {
		c.Ctx.Output.SetStatus(404)
		c.Data["Title"] = "Destination Not Found"
		c.Data["Message"] = "The destination you requested could not be found."
		c.TplName = "errors/404.tpl"
		return
	}

	attractions := services.Attractions.GetAttractions(country.Lat, country.Lng)

	c.Data["Title"] = country.Name
	c.Data["Country"] = country
	c.Data["Attractions"] = attractions
	c.TplName = "destination.tpl"
}
