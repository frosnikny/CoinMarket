package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type InventoryItem struct {
	gorm.Model
	UserID   uuid.UUID `gorm:"type:uuid;not null"`
	ItemType string    `gorm:"not null"`
	Quantity int       `gorm:"not null"`
}
