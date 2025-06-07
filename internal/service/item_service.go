package service

import (
	"github.com/chinaboard/cotify/pkg/model"
	"github.com/chinaboard/cotify/pkg/storage"
)

type ItemService struct {
	storage storage.Storage
}

func NewItemService(storage storage.Storage) *ItemService {
	return &ItemService{
		storage: storage,
	}
}

// StoreItem stores a new item or returns existing one
func (s *ItemService) StoreItem(url, title, itemType, attribute string) (*model.Item, bool, error) {
	return s.storage.StoreItem(url, title, itemType, attribute)
}

// GetItem retrieves an item by URL
func (s *ItemService) GetItem(url string) (*model.Item, error) {
	return s.storage.GetItem(url)
}
