package utils

import (
	"log/slog"

	"github.com/a-h/templ"
	"github.com/gofiber/fiber/v3"
)

func Render(c fiber.Ctx, component templ.Component, status int) error {
	c.Set("Content-Type", "text/html")
	c.Status(status)
	err := component.Render(c.Context(), c.Response().BodyWriter())
	if err != nil {
		slog.Error("Render error", "error", err)
		return err
	}
	return nil
}
