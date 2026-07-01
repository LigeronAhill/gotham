package middleware

import (
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/compress"
)

func Compress() fiber.Handler {
	cfg := compress.Config{
		Level: compress.LevelBestSpeed,
	}
	return compress.New(cfg)
}
