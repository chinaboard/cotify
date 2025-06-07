package app

import (
	"log"

	"github.com/chinaboard/cotify/internal/api"
	"github.com/chinaboard/cotify/internal/config"
	"github.com/chinaboard/cotify/internal/router"
	"github.com/chinaboard/cotify/internal/service"
	"github.com/chinaboard/cotify/pkg/storage"

	"github.com/gin-gonic/gin"
)

type App struct {
	router *gin.Engine
}

func NewApp(cfg *config.Config) *App {
	dsn := cfg.GetDSN()

	storageService, err := storage.NewStorageService(dsn)
	if err != nil {
		log.Fatalf("Failed to initialize storage service: %v", err)
	}

	itemService := service.NewItemService(storageService)
	itemHandler := api.NewItemHandler(itemService)
	r := router.SetupRouter(itemHandler)

	return &App{
		router: r,
	}
}

func (a *App) Run(cfg *config.Config) error {
	return a.router.Run(":" + cfg.ServerPort)
}
