package models

import "gorm.io/gorm"

type InventoryItem struct {
	gorm.Model
	Username string `gorm:"not null"`
	ItemType string `gorm:"not null"`
	Quantity int    `gorm:"not null"`
}
