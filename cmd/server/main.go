package main

import (
	"log"

	"github.com/chinaboard/cotify/app"
	"github.com/chinaboard/cotify/internal/config"
)

func main() {
	cfg := config.LoadFromEnv()
	application := app.NewApp(cfg)
	if err := application.Run(cfg); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
