package handlers

import (
	"CoinMarket/internal/domain/requests"
	"CoinMarket/internal/domain/responses"
	"CoinMarket/internal/usecase"
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
)

type AuthHandler struct {
	service *usecase.AuthService
}

func NewAuthHandler(service *usecase.AuthService) *AuthHandler {
	return &AuthHandler{service: service}
}

func (h *AuthHandler) Auth(c *gin.Context) {
	var request requests.AuthRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	token, err := h.service.Login(request.Username, request.Password)
	if err != nil {
		switch {
		case errors.Is(err, usecase.ErrUserNotFound):
			token, err = h.service.Register(request.Username, request.Password)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to register user"})
				return
			}
		case errors.Is(err, usecase.ErrInvalidPass):
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid password"})
			return
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
			return
		}
	}

	c.JSON(http.StatusOK, responses.AuthResponse{Token: token})
}
