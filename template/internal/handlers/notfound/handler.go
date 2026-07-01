package notfound

import (
	"log/slog"

	"{{ .ModulePath }}/internal/utils"

	"github.com/gofiber/fiber/v3"
)

type handler struct{}

func NewHandler(app fiber.Router) {
	h := &handler{}
	app.All("*", h.get)
}

func (h *handler) get(c fiber.Ctx) error {
	slog.Info("route not found")
	component := view()
	return utils.Render(c, component, fiber.StatusNotFound)
}
