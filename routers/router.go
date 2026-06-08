package routers

import (
	"log"
	"net/http"
	"strings"
	"time"

	"travel-sphere-rangon/controllers"
	api "travel-sphere-rangon/controllers/api"

	beego "github.com/beego/beego/v2/server/web"
	"github.com/beego/beego/v2/server/web/context"
)

func init() {
	beego.Router("/", &controllers.HomeController{})
	beego.Router("/login", &controllers.LoginController{})
	beego.Router("/logout", &controllers.LogoutController{})

	beego.Router("/countries", &controllers.CountryController{})
	beego.Router("/countries/:slug", &controllers.CountryDetailController{})

	beego.Router("/wishlist", &controllers.WishlistController{})
	beego.Router("/dashboard", &controllers.DashboardController{})

	beego.Router("/api/countries", &api.CountriesController{})
	beego.Router("/api/countries/:slug", &api.CountryDetailController{})

	beego.Router("/api/wishlist", &api.WishlistController{})
	beego.Router("/api/wishlist/:id", &api.WishlistItemController{})

	beego.Router("/api/dashboard/summary", &api.DashboardSummaryController{})

	beego.InsertFilter("/*", beego.BeforeRouter, func(ctx *context.Context) {
		ctx.Input.SetData("startTime", time.Now())
	})

	beego.InsertFilter("/*", beego.FinishRouter, func(ctx *context.Context) {
		startTime, ok := ctx.Input.GetData("startTime").(time.Time)
		if !ok {
			return
		}

		log.Printf("%s %s %s", ctx.Request.Method, ctx.Request.URL.Path, time.Since(startTime))
	})

	beego.InsertFilter("/wishlist", beego.BeforeRouter, requireAuth)
	beego.InsertFilter("/dashboard", beego.BeforeRouter, requireAuth)
	beego.InsertFilter("/api/wishlist", beego.BeforeRouter, requireAuth)
	beego.InsertFilter("/api/wishlist/*", beego.BeforeRouter, requireAuth)
	beego.InsertFilter("/api/dashboard/summary", beego.BeforeRouter, requireAuth)
}

func requireAuth(ctx *context.Context) {
	username, ok := ctx.Input.Session("username").(string)
	if ok && strings.TrimSpace(username) != "" {
		return
	}

	if strings.HasPrefix(ctx.Request.URL.Path, "/api/") {
		ctx.Output.SetStatus(http.StatusUnauthorized)
		_ = ctx.Output.JSON(map[string]interface{}{
			"message": "Authentication required.",
			"status":  http.StatusUnauthorized,
		}, false, false)
		return
	}

	ctx.Redirect(http.StatusFound, "/login")
}
