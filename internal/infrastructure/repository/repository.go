package repository

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Repository struct {
	DB *gorm.DB
}

var _ MainRepositoryInterface = (*Repository)(nil)

func New(dsn string) (*Repository, error) {
	db, err := CreateDB(dsn)
	if err != nil {
		return nil, err
	}

	return &Repository{
		DB: db,
	}, nil
}

func CreateDB(dsn string) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	return db, nil
}
