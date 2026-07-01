package middleware

import (
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/static"
)

func Static(isProduction bool) fiber.Handler {
	var cfg static.Config
	if isProduction {
		cfg = static.Config{
			CacheDuration: 10 * time.Second,
			MaxAge:        86400,
			Compress:      true,
			ByteRange:     true,
			Browse:        false,
			Download:      false,
		}
	} else {
		cfg = static.Config{
			CacheDuration: -1 * time.Second,
			MaxAge:        0,
			Compress:      false,
			ByteRange:     false,
			Browse:        true,
			Download:      false,
		}
	}
	return static.New("./public", cfg)
}
