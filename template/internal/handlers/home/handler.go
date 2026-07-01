package home

import (
	"{{ .ModulePath }}/internal/utils"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/csrf"
)

type handler struct{}

func NewHandler(app fiber.Router) {
	h := &handler{}
	app.Get("/", h.get)
}

func (h *handler) get(c fiber.Ctx) error {
	token := csrf.TokenFromContext(c)
	component := view(token)
	return utils.Render(c, component, fiber.StatusOK)
}
