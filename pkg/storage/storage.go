package storage

import (
	"errors"
	"time"

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
	db *gorm.DB
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

	return &StorageService{db: db}, nil
}

// StoreItem stores a new item or returns existing one
func (s *StorageService) StoreItem(url, title, itemType, metadata string) (*model.Item, bool, error) {
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
			return newItem, true, nil
		}
		return nil, false, result.Error
	}

	// Item already exists
	return &item, false, nil
}

// GetItem retrieves an item by URL
func (s *StorageService) GetItem(url string) (*model.Item, error) {
	var item model.Item
	result := s.db.Where("url = ?", url).First(&item)
	if result.Error != nil {
		return nil, result.Error
	}
	return &item, nil
}

// ListItems retrieves all items with optional filtering
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
