package api

import (
	"encoding/json"
	"net/http"
	"travel-sphere-rangon/models"
	"travel-sphere-rangon/services"

	beego "github.com/beego/beego/v2/server/web"
)

type WishlistController struct {
	beego.Controller
}

type wishlistCreateRequest struct {
	CountryName string `json:"country_name"`
	Note        string `json:"note"`
	Status      string `json:"status"`
}

type wishlistUpdateRequest struct {
	Note   string `json:"note"`
	Status string `json:"status"`
}

func currentUser(c *beego.Controller) string {
	username, _ := c.GetSession("username").(string)
	return username
}

func (c *WishlistController) Get() {
	username := currentUser(&c.Controller)
	items := services.Wishlist.List(username)

	c.Data["json"] = items
	c.ServeJSON()
}

// When a POST request comes to the wishlist endpoint,
// read the JSON body,
// validate and create a wishlist item through the service,
// then return either an error response or the created item as JSON.

func (c *WishlistController) Post() {
	username := currentUser(&c.Controller)

	var req wishlistCreateRequest
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &req); err != nil {
		c.Ctx.Output.SetStatus(http.StatusBadRequest)
		c.Data["json"] = models.ErrorResponse{
			Message: "Invalid JSON payload.",
			Status:  http.StatusBadRequest,
		}
		c.ServeJSON()
		return
	}

	item, err := services.Wishlist.Create(username, req.CountryName, req.Note, req.Status)
	if err != nil {
		c.Ctx.Output.SetStatus(http.StatusBadRequest)
		c.Data["json"] = models.ErrorResponse{
			Message: err.Error(),
			Status:  http.StatusBadRequest,
		}
		c.ServeJSON()
		return
	}

	c.Ctx.Output.SetStatus(http.StatusCreated)
	c.Data["json"] = item
	c.ServeJSON()
}

type WishlistItemController struct {
	beego.Controller
}

func (c *WishlistItemController) Put() {
	username := currentUser(&c.Controller)
	id := c.Ctx.Input.Param(":id")

	var req wishlistUpdateRequest
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &req); err != nil {
		c.Ctx.Output.SetStatus(http.StatusBadRequest)
		c.Data["json"] = models.ErrorResponse{
			Message: "Invalid JSON payload.",
			Status:  http.StatusBadRequest,
		}
		c.ServeJSON()
		return
	}

	item, err := services.Wishlist.Update(username, id, req.Note, req.Status)
	if err != nil {
		status := http.StatusBadRequest
		if err.Error() == "wishlist item not found" {
			status = http.StatusNotFound
		}

		c.Ctx.Output.SetStatus(status)
		c.Data["json"] = models.ErrorResponse{
			Message: err.Error(),
			Status:  status,
		}
		c.ServeJSON()
		return
	}

	c.Data["json"] = item
	c.ServeJSON()
}

func (c *WishlistItemController) Delete() {
	username := currentUser(&c.Controller)
	id := c.Ctx.Input.Param(":id")

	if err := services.Wishlist.Delete(username, id); err != nil {
		c.Ctx.Output.SetStatus(http.StatusNotFound)
		c.Data["json"] = models.ErrorResponse{
			Message: err.Error(),
			Status:  http.StatusNotFound,
		}
		c.ServeJSON()
		return
	}

	c.Data["json"] = models.SuccessResponse{
		Message: "Wishlist item deleted.",
	}
	c.ServeJSON()
}
