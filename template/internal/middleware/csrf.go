package middleware

import (
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/extractors"
	"github.com/gofiber/fiber/v3/middleware/csrf"
	"github.com/gofiber/fiber/v3/middleware/session"
)

func CSRF(isProduction bool, store *session.Store) fiber.Handler {
	cookieName := "csrf_"
	if isProduction {
		cookieName = "__Host-csrf_"
	}
	cfg := csrf.Config{
		Session:           store,
		CookieName:        cookieName,
		CookieSameSite:    "Lax",
		IdleTimeout:       30 * time.Minute,
		CookieSecure:      isProduction,
		CookieHTTPOnly:    true,
		CookieSessionOnly: true,
		Extractor:         extractors.Chain(extractors.FromHeader("X-CSRF-Token"), extractors.FromForm("_csrf")),
	}
	return csrf.New(cfg)
}
