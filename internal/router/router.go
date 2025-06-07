package router

import (
	"github.com/chinaboard/cotify/internal/api"

	"github.com/gin-gonic/gin"
)

// SetupRouter configures all the routes for the application
func SetupRouter(itemHandler *api.ItemHandler) *gin.Engine {
	r := gin.Default()

	// API routes
	apiGroup := r.Group("/api")
	{
		apiGroup.POST("/items", itemHandler.StoreItem)
	}

	return r
}
