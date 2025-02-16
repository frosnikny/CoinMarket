package middleware

import (
	"CoinMarket/internal/infrastructure/repository"
	"CoinMarket/internal/usecase"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

func AuthMiddleware(authService *usecase.AuthService, userRepo *repository.UserRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")

		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing or invalid Authorization header"})
			c.Abort()
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		claims, err := authService.ValidateToken(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			c.Abort()
			return
		}

		user, err := userRepo.GetUserByUsername(claims.Username)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
			c.Abort()
			return
		}

		c.Set("userID", user.ID)
		c.Set("username", user.Username)

		c.Next()
	}
}
