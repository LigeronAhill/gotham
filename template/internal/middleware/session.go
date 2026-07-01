package middleware

import (
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/extractors"
	"github.com/gofiber/fiber/v3/middleware/session"
)

func Session(isProduction bool, storage fiber.Storage) (fiber.Handler, *session.Store) {
	cookieName := "session_id"
	if isProduction {
		cookieName = "__Host-session_id"
	}
	cfg := session.Config{
		Storage:         storage,
		CookieSameSite:  "Lax",
		Extractor:       extractors.FromCookie(cookieName),
		IdleTimeout:     30 * time.Minute,
		AbsoluteTimeout: 24 * time.Hour,
		CookieSecure:    isProduction,
		CookieHTTPOnly:  true,
	}
	return session.NewWithStore(cfg)
}
