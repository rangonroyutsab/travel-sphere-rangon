package controllers

type WishlistController struct {
	BaseController
}

func (c *WishlistController) Get() {
	c.Data["Title"] = "My Wishlist"
	c.TplName = "wishlist.tpl"
}
