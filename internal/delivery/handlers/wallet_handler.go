package handlers

import (
	"CoinMarket/internal/domain/requests"
	"CoinMarket/internal/domain/responses"
	"CoinMarket/internal/usecase"
	"errors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
)

type WalletHandler struct {
	service usecase.WalletServiceInterface
}

func NewWalletHandler(service usecase.WalletServiceInterface) *WalletHandler {
	return &WalletHandler{service: service}
}

func (h *WalletHandler) GetUserInfo(c *gin.Context) {
	username, exists := c.Get("username")
	if !exists {
		c.JSON(http.StatusUnauthorized, responses.ErrorResponse{Errors: "Unauthorized"})
		return
	}

	usernameStr, ok := username.(string)
	if !ok {
		c.JSON(http.StatusInternalServerError, responses.ErrorResponse{Errors: "Invalid user ID"})
		return
	}

	info, err := h.service.GetUserInfo(usernameStr)
	if err != nil {
		c.JSON(http.StatusInternalServerError, responses.ErrorResponse{Errors: "Failed to fetch user info"})
		return
	}

	c.JSON(http.StatusOK, info)
}

func (h *WalletHandler) SendCoins(c *gin.Context) {
	fromUsername, exists := c.Get("username")
	if !exists {
		c.JSON(http.StatusUnauthorized, responses.ErrorResponse{Errors: "Unauthorized"})
		return
	}
	fromUsernameStr, ok := fromUsername.(string)
	if !ok {
		c.JSON(http.StatusInternalServerError, responses.ErrorResponse{Errors: "Invalid user ID"})
		return
	}

	var req requests.SendCoinRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, responses.ErrorResponse{Errors: "Invalid request"})
		return
	}

	toUser := req.ToUser
	err := h.service.SendCoins(fromUsernameStr, toUser, req.Amount)
	if err != nil {
		if errors.Is(err, usecase.ErrInsufficientFunds) {
			c.JSON(http.StatusBadRequest, responses.ErrorResponse{Errors: "Not enough coins"})
		} else if errors.Is(err, usecase.ErrNegativeAmount) {
			c.JSON(http.StatusBadRequest, responses.ErrorResponse{Errors: "Amount must be positive"})
		} else if errors.Is(err, usecase.ErrSelfTransfer) {
			c.JSON(http.StatusBadRequest, responses.ErrorResponse{Errors: "Cannot send coins to yourself"})
		} else if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusBadRequest, responses.ErrorResponse{Errors: "Recipient does not exist"})
		} else {
			c.JSON(http.StatusInternalServerError, responses.ErrorResponse{Errors: "Transaction failed"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Coins sent successfully"})
}

func (h *WalletHandler) BuyItem(c *gin.Context) {
	username, exists := c.Get("username")
	if !exists {
		c.JSON(http.StatusUnauthorized, responses.ErrorResponse{Errors: "Unauthorized"})
		return
	}
	usernameStr, ok := username.(string)
	if !ok {
		c.JSON(http.StatusInternalServerError, responses.ErrorResponse{Errors: "Invalid user ID"})
		return
	}

	itemName := c.Param("item")
	err := h.service.BuyItem(usernameStr, itemName)
	if err != nil {
		if errors.Is(err, usecase.ErrItemNotFound) {
			c.JSON(http.StatusNotFound, responses.ErrorResponse{Errors: "Item not found"})
		} else if errors.Is(err, usecase.ErrNotEnoughCoins) {
			c.JSON(http.StatusBadRequest, responses.ErrorResponse{Errors: "Not enough coins"})
		} else {
			c.JSON(http.StatusInternalServerError, responses.ErrorResponse{Errors: "Failed to buy item"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Item purchased successfully"})
}
