package storage

import (
	"errors"
	"time"

	"github.com/chinaboard/cotify/pkg/cache"
	"github.com/chinaboard/cotify/pkg/model"
	gorm_logrus "github.com/onrik/gorm-logrus"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// Storage defines the interface for storage operations
type Storage interface {
	StoreItem(url, title, itemType, metadata string) (*model.Item, bool, error)
	GetItem(url string) (*model.Item, error)
	ListItems(itemType string, startTime, endTime *time.Time) ([]model.Item, error)
}

// StorageService implements the Storage interface
type StorageService struct {
	db    *gorm.DB
	cache *cache.MemoryCache
}

// NewStorageService creates a new storage service instance
func NewStorageService(dsn string) (Storage, error) {
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: gorm_logrus.New(),
	})
	if err != nil {
		return nil, err
	}

	// Auto migrate the schema
	err = db.AutoMigrate(&model.Item{})
	if err != nil {
		return nil, err
	}

	service := &StorageService{
		db:    db,
		cache: cache.NewMemoryCache(24 * time.Hour), // 24 hours TTL
	}

	return service, nil
}

// StoreItem stores a new item or returns existing one
func (s *StorageService) StoreItem(url, title, itemType, metadata string) (*model.Item, bool, error) {
	// Check cache first
	if cachedValue, found := s.cache.Get(url); found {
		if item, ok := cachedValue.(*model.Item); ok {
			return item, false, nil
		}
	}

	var item model.Item
	result := s.db.Where("url = ?", url).First(&item)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			// Create new item
			newItem := &model.Item{
				Title:    title,
				Url:      url,
				Type:     itemType,
				Metadata: metadata,
			}
			if err := s.db.Create(newItem).Error; err != nil {
				return nil, false, err
			}
			// Cache the new item
			s.cache.Set(url, newItem)
			return newItem, true, nil
		}
		return nil, false, result.Error
	}

	// Item already exists, cache it
	s.cache.Set(url, &item)
	return &item, false, nil
}

// GetItem retrieves an item by URL
func (s *StorageService) GetItem(url string) (*model.Item, error) {
	// Check cache first
	if cachedValue, found := s.cache.Get(url); found {
		if item, ok := cachedValue.(*model.Item); ok {
			return item, nil
		}
	}

	var item model.Item
	result := s.db.Where("url = ?", url).First(&item)
	if result.Error != nil {
		return nil, result.Error
	}

	// Cache the item
	s.cache.Set(url, &item)
	return &item, nil
}

// ListItems retrieves all items with optional filtering
// Note: ListItems is not cached as it's more complex to cache filtered results
func (s *StorageService) ListItems(itemType string, startTime, endTime *time.Time) ([]model.Item, error) {
	var items []model.Item
	query := s.db.Model(&model.Item{})

	if itemType != "" {
		query = query.Where("type = ?", itemType)
	}

	if startTime != nil {
		query = query.Where("created_at >= ?", startTime)
	}
	if endTime != nil {
		query = query.Where("created_at <= ?", endTime)
	}

	err := query.Find(&items).Error
	return items, err
}
