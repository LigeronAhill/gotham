package main

import (
	"context"

	"{{ .ModulePath }}/internal/config"
	"{{ .ModulePath }}/internal/logger"
	"{{ .ModulePath }}/internal/server"
)

func main() {
	ctx := context.Background()
	settings := config.Init()
	logger := logger.Init(settings)
	app := server.New(settings, logger)
	if err := app.Run(ctx); err != nil {
		panic(err)
	}
}
