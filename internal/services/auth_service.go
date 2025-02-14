package services

import (
	"CoinMarket/internal/models"
	"CoinMarket/internal/repository"
	"errors"
	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
	"time"
)

type AuthService struct {
	repo   *repository.UserRepository
	jwtKey []byte
}

func NewAuthService(repo *repository.UserRepository, jwtKey string) *AuthService {
	return &AuthService{
		repo:   repo,
		jwtKey: []byte(jwtKey),
	}
}

type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

// TODO: использовать текст ошибки ?
var (
	ErrUserNotFound = errors.New("user not found")
	ErrInvalidPass  = errors.New("invalid password")
)

func (s *AuthService) Register(username, password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	user := models.User{Username: username, Password: string(hashedPassword)}
	err = s.repo.CreateUser(&user)
	if err != nil {
		return "", err
	}

	return s.GenerateToken(username)
}

func (s *AuthService) Login(username, password string) (string, error) {
	user, err := s.repo.GetUserByUsername(username)
	if err != nil {
		return "", ErrUserNotFound
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return "", ErrInvalidPass
	}

	return s.GenerateToken(username)
}

func (s *AuthService) GenerateToken(username string) (string, error) {
	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &Claims{
		Username: username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
			Issuer:    "CoinMarket",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.jwtKey)
}

func (s *AuthService) ValidateToken(tokenString string) (*Claims, error) {
	claims := &Claims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return s.jwtKey, nil
	})

	if err != nil || !token.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}
