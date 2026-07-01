package middleware

import (
	"log/slog"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/recover"
	"github.com/gofiber/fiber/v3/middleware/requestid"
	slogfiber "github.com/samber/slog-fiber"
)

func Setup(app *fiber.App, origins []string, isProduction bool, storage fiber.Storage, logger *slog.Logger) {
	loggerConfig := slogfiber.Config{
		WithRequestID:    true,
		DefaultLevel:     slog.LevelInfo,
		ClientErrorLevel: slog.LevelWarn,
		ServerErrorLevel: slog.LevelError,
	}

	cors := Cors(origins)

	helmet := Helmet(isProduction)

	compress := Compress()

	favicon := Favicon()

	static := Static(isProduction)

	limiter := Limiter(storage)

	cache := Cache(storage)

	sess, store := Session(isProduction, storage)

	csrf := CSRF(isProduction, store)

	app.Use(recover.New())

	app.Use(helmet)

	app.Use(compress)

	app.Use(cors)

	app.Use(requestid.New())

	app.Use(slogfiber.NewWithConfig(logger, loggerConfig))

	app.Use(limiter)

	app.Use(cache)

	app.Use(sess)

	app.Use(csrf)

	app.Use(favicon)

	app.Use("/public", static)
}
