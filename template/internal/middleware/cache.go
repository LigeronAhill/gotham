package middleware

import (
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cache"
)

func Cache(storage fiber.Storage) fiber.Handler {
	cfg := cache.Config{
		Storage: storage,
		Next: func(c fiber.Ctx) bool {
			return c.Get("HX-Request") == "true" && c.Get("HX-Boosted") == "true"
		},
		KeyGenerator: func(c fiber.Ctx) string {
			return c.Method() + ":" + c.Path() + "?" + c.RequestCtx().QueryArgs().String() + "|hx=" + c.Get("HX-Request") + "|target=" + c.Get("HX-Target")
		},
		Methods: []string{
			fiber.MethodGet,
			fiber.MethodHead,
		},
		Expiration:           2 * time.Minute,
		StoreResponseHeaders: true,
	}
	return cache.New(cfg)
}
