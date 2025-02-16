package integration

import (
	"CoinMarket/internal/delivery/routes"
	"CoinMarket/internal/infrastructure/db/dsn"
	"CoinMarket/internal/infrastructure/db/migrations"
	"CoinMarket/internal/infrastructure/repository"
	"bytes"
	"encoding/json"
	"gorm.io/gorm"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"CoinMarket/internal/app"
	"CoinMarket/internal/config"
	"github.com/gin-gonic/gin"
	"github.com/ory/dockertest/v3"
	"github.com/stretchr/testify/assert"
)

var testRouter *gin.Engine
var testApp *app.Application
var testDB *gorm.DB

func TestMain(m *testing.M) {
	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Fatalf("Could not connect to Docker: %s", err)
	}

	resource, err := pool.Run("postgres", "15", []string{
		"POSTGRES_USER=test",
		"POSTGRES_PASSWORD=test",
		"POSTGRES_DB=testdb",
	})
	if err != nil {
		log.Fatalf("Could not start resource: %s", err)
	}

	// Ждем поднятия БД
	time.Sleep(5 * time.Second)

	cfg := &config.Config{
		PostgresHost:     "localhost",
		PostgresPort:     resource.GetPort("5432/tcp"),
		PostgresUsername: "test",
		PostgresPassword: "test",
		PostgresDatabase: "testdb",
		JwtKey:           "testsecret",
	}

	testDB, err = repository.CreateDB(dsn.FromCfg(cfg))
	if err != nil {
		log.Fatalf("Ошибка подключения к тестовой базе данных: %v", err)
	}

	err = migrations.RunMigrations(testDB)
	if err != nil {
		log.Fatalf("Ошибка миграции: %v", err)
	}

	testApp = app.New(cfg)
	testRouter = gin.Default()
	routes.SetupRoutes(testRouter, testApp)

	code := m.Run()

	pool.Purge(resource)

	os.Exit(code)
}

func TestMerchPurchase(t *testing.T) {
	// 1. Авторизация: получаем токен
	authPayload := `{"username": "testuser", "password": "testpass"}`
	req, _ := http.NewRequest("POST", "/api/auth", bytes.NewBuffer([]byte(authPayload)))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	testRouter.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	var authResp map[string]string
	json.Unmarshal(w.Body.Bytes(), &authResp)
	token := authResp["token"]
	assert.NotEmpty(t, token)

	// 2. Проверяем баланс перед покупкой
	req, _ = http.NewRequest("GET", "/api/info", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w = httptest.NewRecorder()
	testRouter.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	var infoResp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &infoResp)

	coinsBefore := int(infoResp["coins"].(float64))
	assert.True(t, coinsBefore >= 80, "User must have at least 80 coins to buy a t-shirt")

	// 3. Покупаем товар "t-shirt"
	itemID := "t-shirt"
	req, _ = http.NewRequest("GET", "/api/buy/"+itemID, nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w = httptest.NewRecorder()
	testRouter.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	// 4. Проверяем баланс и инвентарь после покупки
	req, _ = http.NewRequest("GET", "/api/info", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w = httptest.NewRecorder()
	testRouter.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
	json.Unmarshal(w.Body.Bytes(), &infoResp)

	coinsAfter := int(infoResp["coins"].(float64))
	assert.True(t, coinsAfter == coinsBefore-80, "Balance should decrease by 80 coins")

	// Проверяем, что "t-shirt" добавился в инвентарь
	newItems := infoResp["inventory"].([]interface{})
	found := false
	for _, i := range newItems {
		if i.(map[string]interface{})["ItemType"].(string) == itemID {
			found = true
			break
		}
	}
	assert.True(t, found, "Purchased item (t-shirt) should be in inventory")
}
