package controllers

import "travel-sphere-rangon/services"

type DashboardController struct {
	BaseController
}

func (c *DashboardController) Get() {
	username := c.CurrentUsername()
	summary := services.Dashboard.GetSummary(username)

	c.Data["Title"] = "Dashboard"
	c.Data["TotalSaved"] = summary.TotalSaved
	c.Data["Planned"] = summary.Planned
	c.Data["Visited"] = summary.Visited
	c.Data["SavedDestinations"] = services.Dashboard.GetSavedDestinations(username)

	c.TplName = "dashboard.tpl"
}
