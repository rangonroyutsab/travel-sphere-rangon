package api

import beego "github.com/beego/beego/v2/server/web"

type WishlistController struct {
	beego.Controller
}

func (c *WishlistController) Get() {
	c.Data["json"] = []string{}
	c.ServeJSON()
}

func (c *WishlistController) Post() {
	c.Data["json"] = map[string]string{
		"message": "wishlist create",
	}
	c.ServeJSON()
}

type WishlistItemController struct {
	beego.Controller
}

func (c *WishlistItemController) Put() {
	c.Data["json"] = map[string]string{
		"message": "wishlist update",
	}
	c.ServeJSON()
}

func (c *WishlistItemController) Delete() {
	c.Data["json"] = map[string]string{
		"message": "wishlist delete",
	}
	c.ServeJSON()
}
