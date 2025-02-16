package unit

import (
	"CoinMarket/internal/delivery/handlers"
	"CoinMarket/internal/tests/mocks/servicesmocks"
	"CoinMarket/internal/usecase"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandlerGetUserInfo_Success(t *testing.T) {
	mockService := new(servicesmocks.WalletServiceMock)
	handler := handlers.NewWalletHandler(mockService)
	username := "test_user"
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Set("username", username)

	mockService.On("GetUserInfo", username).Return(&usecase.InfoResponse{Coins: 100}, nil)

	handler.GetUserInfo(c)

	assert.Equal(t, http.StatusOK, w.Code)
	mockService.AssertExpectations(t)
}

func TestHandlerBuyItem_Success(t *testing.T) {
	mockService := new(servicesmocks.WalletServiceMock)
	handler := handlers.NewWalletHandler(mockService)
	username := "test_user"
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Set("username", username)
	c.Params = []gin.Param{{Key: "item", Value: "sword"}}

	mockService.On("BuyItem", username, "sword").Return(nil)

	handler.BuyItem(c)

	assert.Equal(t, http.StatusOK, w.Code)
	mockService.AssertExpectations(t)
}
