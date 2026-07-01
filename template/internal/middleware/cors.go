package middleware

import (
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
)

func Cors(origins []string) fiber.Handler {
	cfg := cors.Config{
		AllowOrigins: origins,
		AllowMethods: []string{
			fiber.MethodGet,
			fiber.MethodPost,
			fiber.MethodPut,
			fiber.MethodPatch,
			fiber.MethodDelete,
			fiber.MethodOptions,
		},
		AllowHeaders: []string{
			"Origin",
			"Content-Type",
			"Accept",
			"Authorization",
			"X-Requested-With",
			"HX-Request",
			"HX-Target",
			"HX-Trigger",
			"HX-Trigger-Name",
			"HX-History-Restore-Request",
			"HX-Current-URL",
			"HX-Prompt",
			"HX-Boosted",
		},
		ExposeHeaders: []string{
			"HX-Trigger",
			"HX-Location",
			"HX-Refresh",
			"HX-Redirect",
		},
		MaxAge:           86400,
		AllowCredentials: true,
	}
	return cors.New(cfg)
}
