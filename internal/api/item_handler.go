package api

import (
	"net/http"

	"github.com/chinaboard/cotify/internal/service"
	"github.com/chinaboard/cotify/pkg/model"

	"github.com/gin-gonic/gin"
)

type ItemHandler struct {
	itemService *service.ItemService
}

func NewItemHandler(itemService *service.ItemService) *ItemHandler {
	return &ItemHandler{
		itemService: itemService,
	}
}

type StoreItemRequest struct {
	Url      string `json:"url" binding:"required"`
	Title    string `json:"title" binding:"required"`
	Type     string `json:"type" binding:"required"`
	Metadata string `json:"metadata"`
}

type StoreItemResponse struct {
	Item  *model.Item `json:"item"`
	IsNew bool        `json:"is_new"`
}

func (h *ItemHandler) StoreItem(c *gin.Context) {
	var req StoreItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	item, isNew, err := h.itemService.StoreItem(req.Url, req.Title, req.Type, req.Metadata)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, StoreItemResponse{
		Item:  item,
		IsNew: isNew,
	})
}
