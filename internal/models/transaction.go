package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Transaction struct {
	gorm.Model
	FromUserID uuid.UUID `gorm:"type:uuid;not null"`
	ToUserID   uuid.UUID `gorm:"type:uuid;not null"`
	Amount     int       `gorm:"not null"`
}
