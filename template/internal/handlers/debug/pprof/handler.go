package pprof

import (
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/basicauth"
	"github.com/gofiber/fiber/v3/middleware/pprof"
)

func NewHandler(app *fiber.App, adminPassword string) {
	g := app.Group("/debug/pprof")
	cfg := basicauth.Config{
		Authorizer: func(username, password string, c fiber.Ctx) bool {
			return username == "admin" && password == adminPassword
		},
	}
	ba := basicauth.New(cfg)
	g.Use(ba)
	pp := pprof.New()
	g.Use(pp)
}
