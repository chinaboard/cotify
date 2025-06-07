package main

import (
	"log"

	"github.com/chinaboard/cotify/app"
	"github.com/chinaboard/cotify/internal/config"
)

func main() {
	// cfg := config.NewConfig("localhost", "3306", "cotify_test", "cotify_test", "cotify_test", "8888")
	cfg := config.LoadFromEnv()
	application := app.NewApp(cfg)
	if err := application.Run(cfg); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
