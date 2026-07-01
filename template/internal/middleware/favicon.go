package middleware

import (
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/favicon"
)

func Favicon() fiber.Handler {
	cfg := favicon.Config{
		File: "./public/icons/favicon.ico",
		URL:  "/favicon.ico",
	}
	return favicon.New(cfg)
}
