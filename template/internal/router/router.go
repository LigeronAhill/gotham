package router

import (
	"log/slog"

	"{{ .ModulePath }}/internal/handlers/debug/pprof"
	"{{ .ModulePath }}/internal/handlers/home"
	"{{ .ModulePath }}/internal/handlers/notfound"

	"github.com/gofiber/fiber/v3"
	"github.com/jackc/pgx/v5/pgxpool"
)

func New(pool *pgxpool.Pool, adminPassword string, logger *slog.Logger, app *fiber.App) *fiber.App {
	home.NewHandler(app)

	pprof.NewHandler(app, adminPassword)

	notfound.NewHandler(app)
	return app
}
