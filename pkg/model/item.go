package model

import (
	"gorm.io/gorm"
)

// Item represents a stored HTTP link with additional metadata
type Item struct {
	Title    string `gorm:"type:varchar(2048);not null"`
	Url      string `gorm:"type:varchar(2048);not null;uniqueIndex"`
	Type     string `gorm:"type:varchar(50);not null"`
	Metadata string `gorm:"type:text"`
	gorm.Model
}

// TableName specifies the table name for the Item model
func (Item) TableName() string {
	return "items"
}
