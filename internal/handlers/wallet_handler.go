package handlers

import (
	"CoinMarket/internal/services"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
)

type WalletHandler struct {
	service *services.WalletService
}

func NewWalletHandler(service *services.WalletService) *WalletHandler {
	return &WalletHandler{service: service}
}

func (h *WalletHandler) GetUserInfo(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	userUUID, ok := userID.(uuid.UUID)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID"})
		return
	}

	info, err := h.service.GetUserInfo(userUUID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch user info"})
		return
	}

	c.JSON(http.StatusOK, info)
}

type SendCoinRequest struct {
	ToUser string `json:"toUser"`
	Amount int    `json:"amount"`
}

func (h *WalletHandler) SendCoins(c *gin.Context) {
	fromUserID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	fromUUID, ok := fromUserID.(uuid.UUID)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID"})
		return
	}

	var req SendCoinRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	toUser, err := h.service.UserRepo.GetUserByUsername(req.ToUser)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Recipient not found"})
		return
	}

	err = h.service.SendCoins(fromUUID, toUser.ID, req.Amount)
	if err != nil {
		if errors.Is(err, services.ErrInsufficientFunds) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Not enough coins"})
		} else if errors.Is(err, services.ErrNegativeAmount) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Amount must be positive"})
		} else if errors.Is(err, services.ErrSelfTransfer) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Cannot send coins to yourself"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Transaction failed"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Coins sent successfully"})
}
