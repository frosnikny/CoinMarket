package handlers

import (
	"CoinMarket/internal/services"
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
)

type AuthHandler struct {
	service *services.AuthService
}

func NewAuthHandler(service *services.AuthService) *AuthHandler {
	return &AuthHandler{service: service}
}

func (h *AuthHandler) Auth(c *gin.Context) {
	var request struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	token, err := h.service.Login(request.Username, request.Password)
	if err != nil {
		if errors.Is(err, services.ErrUserNotFound) {
			token, err = h.service.Register(request.Username, request.Password)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to register user"})
				return
			}
		} else if errors.Is(err, services.ErrInvalidPass) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid password"})
			return
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}
