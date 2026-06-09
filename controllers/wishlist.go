package controllers

import "travel-sphere-rangon/services"

type WishlistController struct {
	BaseController
}

func (c *WishlistController) Get() {
	c.Data["Title"] = "My Wishlist"
	c.Data["WishlistItems"] = services.Wishlist.List(c.CurrentUsername())
	c.TplName = "wishlist.tpl"
}
