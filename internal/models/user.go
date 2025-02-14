package models

import (
	"github.com/google/uuid"
	_ "gorm.io/gorm"
)

type User struct {
	ID       uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	Username string    `gorm:"unique;not null"`
	Password string    `gorm:"not null"`
	Coins    int       `gorm:"default:1000"`
}
