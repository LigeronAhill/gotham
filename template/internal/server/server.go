package server

import (
	"context"
	"errors"
	"log/slog"
	"net"
	"strconv"
	"strings"
	"time"

	"{{ .ModulePath }}/internal/database"
	"{{ .ModulePath }}/internal/middleware"
	"{{ .ModulePath }}/internal/router"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/storage/postgres/v3"

	"github.com/spf13/viper"
)

type Server struct {
	host          string
	port          string
	isProduction  bool
	dbURL         string
	adminPassword string
	logger        *slog.Logger
}

func New(settings *viper.Viper, logger *slog.Logger) *Server {
	host := strings.ToLower(settings.GetString("host"))
	port := strconv.Itoa(settings.GetInt("port"))
	env := strings.ToLower(settings.GetString("env"))
	isProduction := env == "production"
	dbURL := settings.GetString("db_url")
	if dbURL == "" {
		panic(errors.New("database url environment variable must be set"))
	}
	adminPassword := settings.GetString("admin_password")
	return &Server{
		host,
		port,
		isProduction,
		dbURL,
		adminPassword,
		logger,
	}
}

func (s *Server) Run(ctx context.Context) error {
	cfg := fiber.Config{
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  30 * time.Second,
		AppName:      "Gotham",
	}
	app := fiber.New(cfg)

	addr := net.JoinHostPort(s.host, s.port)
	var origins []string
	mainOrigin := "http://" + addr
	if s.isProduction {
		mainOrigin = "https://" + addr
	} else {
		origins = append(origins, "http://localhost:7331")
	}
	origins = append(origins, mainOrigin)

	pool, err := database.GetPool(ctx, s.dbURL)
	if err != nil {
		return err
	}
	storage := postgres.New(postgres.Config{
		DB:         pool,
		Table:      "sessions",
		Reset:      false,
		GCInterval: 10 * time.Second,
	})

	middleware.Setup(app, origins, s.isProduction, storage, s.logger)

	router := router.New(pool, s.adminPassword, s.logger, app)

	return router.Listen(addr)
}
