package repository

import (
	"CoinMarket/internal/domain/models"
	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

var _ UserRepositoryInterface = (*UserRepository)(nil)

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) CreateUser(user *models.User) error {
	return r.db.Create(user).Error
}

func (r *UserRepository) GetUserByUsername(username string) (*models.User, error) {
	var user models.User
	err := r.db.Where("username = ?", username).First(&user).Error
	return &user, err
}
