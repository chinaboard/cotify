package model

import (
	"gorm.io/gorm"
)

// Item represents a stored HTTP link with additional metadata
type Item struct {
	gorm.Model
	Title     string `gorm:"type:varchar(2048);not null"`
	URL       string `gorm:"type:varchar(2048);not null;uniqueIndex"`
	Type      string `gorm:"type:varchar(50);not null"`
	Attribute string `gorm:"type:text"`
}

// TableName specifies the table name for the Item model
func (Item) TableName() string {
	return "items"
}
