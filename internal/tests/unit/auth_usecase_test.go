package unit

import (
	"CoinMarket/internal/domain/models"
	"CoinMarket/internal/tests/mocks/repository_mocks"
	"CoinMarket/internal/usecase"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
	"testing"
)

func TestLogin(t *testing.T) {
	mockUserRepo := new(repository_mocks.UserRepositoryMock)
	jwtKey := "test_secret"
	service := usecase.NewAuthService(mockUserRepo, jwtKey)

	username := "testuser"
	password := "password123"
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	mockUserRepo.On("GetUserByUsername", username).
		Return(&models.User{Username: username, Password: string(hashedPassword)}, nil)

	token, err := service.Login(username, password)
	assert.NoError(t, err)
	assert.NotEmpty(t, token)

	mockUserRepo.AssertExpectations(t)
}

func TestGenerateToken(t *testing.T) {
	mockUserRepo := new(repository_mocks.UserRepositoryMock)
	jwtKey := "test_secret"
	service := usecase.NewAuthService(mockUserRepo, jwtKey)

	username := "test"
	token, err := service.GenerateToken(username)
	assert.NoError(t, err)
	assert.NotEmpty(t, token)
}

func TestValidateToken(t *testing.T) {
	mockUserRepo := new(repository_mocks.UserRepositoryMock)
	jwtKey := "test_secret"
	service := usecase.NewAuthService(mockUserRepo, jwtKey)

	username := "test"
	token, _ := service.GenerateToken(username)
	claims, err := service.ValidateToken(token)

	assert.NoError(t, err)
	assert.Equal(t, username, claims.Username)
}
