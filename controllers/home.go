package controllers

import "travel-sphere-rangon/services"

type HomeController struct {
	BaseController
}

func (c *HomeController) Get() {
	c.Data["Title"] = "Home"

	featuredCountries, err := services.Countries.GetFeaturedCountries()
	if err != nil {
		c.Data["CountryError"] = "Featured destinations are temporarily unavailable."
		featuredCountries = nil
	}

	c.Data["FeaturedCountries"] = featuredCountries
	c.TplName = "home.tpl"
}
