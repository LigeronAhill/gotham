package middleware

import (
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/limiter"
)

func Limiter(storage fiber.Storage) fiber.Handler {
	cfg := limiter.Config{
		Storage:    storage,
		Max:        100,
		Expiration: 1 * time.Minute,
	}
	return limiter.New(cfg)
}
