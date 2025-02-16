package unit

import (
	"CoinMarket/internal/domain/models"
	"CoinMarket/internal/tests/mocks/repositorymocks"
	"CoinMarket/internal/usecase"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
	"testing"
)

func TestLogin(t *testing.T) {
	mockUserRepo := new(repositorymocks.UserRepositoryMock)
	jwtKey := "test_secret"
	service := usecase.NewAuthService(mockUserRepo, jwtKey)

	username := "testuser"
	password := "password123"
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	assert.NoError(t, err)

	mockUserRepo.On("GetUserByUsername", username).
		Return(&models.User{Username: username, Password: string(hashedPassword)}, nil)

	token, err := service.Login(username, password)
	assert.NoError(t, err)
	assert.NotEmpty(t, token)

	mockUserRepo.AssertExpectations(t)
}

func TestGenerateToken(t *testing.T) {
	mockUserRepo := new(repositorymocks.UserRepositoryMock)
	jwtKey := "test_secret"
	service := usecase.NewAuthService(mockUserRepo, jwtKey)

	username := "test"
	token, err := service.GenerateToken(username)
	assert.NoError(t, err)
	assert.NotEmpty(t, token)
}

func TestValidateToken(t *testing.T) {
	mockUserRepo := new(repositorymocks.UserRepositoryMock)
	jwtKey := "test_secret"
	service := usecase.NewAuthService(mockUserRepo, jwtKey)

	username := "test"
	token, err := service.GenerateToken(username)
	assert.NoError(t, err)

	claims, err := service.ValidateToken(token)

	assert.NoError(t, err)
	assert.Equal(t, username, claims.Username)
}
