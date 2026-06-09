package api

import (
	"net/http"

	"travel-sphere-rangon/models"
	"travel-sphere-rangon/services"

	beego "github.com/beego/beego/v2/server/web"
)

type CountriesController struct {
	beego.Controller
}

func (c *CountriesController) Get() {
	search := c.GetString("search")
	region := c.GetString("region")

	countries, err := services.Countries.SearchCountries(search, region)
	if err != nil {
		c.Ctx.Output.SetStatus(http.StatusInternalServerError)
		c.Data["json"] = models.ErrorResponse{
			Message: "Could not load countries.",
			Status:  http.StatusInternalServerError,
		}
		c.ServeJSON()
		return
	}
	c.Data["json"] = countries
	c.ServeJSON()
}

type CountryDetailController struct {
	beego.Controller
}

func (c *CountryDetailController) Get() {
	slug := c.Ctx.Input.Param(":slug")

	country, err := services.Countries.GetCountryBySlug(slug)
	if err != nil {
		c.Ctx.Output.SetStatus(http.StatusNotFound)
		c.Data["json"] = models.ErrorResponse{
			Message: "Country not found.",
			Status:  http.StatusNotFound,
		}
		c.ServeJSON()
		return
	}

	c.Data["json"] = country
	c.ServeJSON()
}
