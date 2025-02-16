package integration

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func authenticateUser(username, password string) (string, int) {
	authPayload := map[string]string{
		"username": username,
		"password": password,
	}
	body, _ := json.Marshal(authPayload)

	req, _ := http.NewRequest("POST", "/api/auth", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	testRouter.ServeHTTP(w, req)

	var authResp map[string]string
	json.Unmarshal(w.Body.Bytes(), &authResp)

	return authResp["token"], w.Code
}

func TestAuth(t *testing.T) {
	// 1. Регистрация нового пользователя
	username := "newUser"
	password := "securePass"

	token, status := authenticateUser(username, password)
	assert.Equal(t, http.StatusOK, status, "Registration should succeed")
	assert.NotEmpty(t, token, "JWT token is required on success")

	// 2. Успешная аутентификация существующего
	token, status = authenticateUser(username, password)
	assert.Equal(t, http.StatusOK, status, "Authentication should succeed")
	assert.NotEmpty(t, token, "JWT token is required on success")

	// 3. Ошибка при вводе неправильного пароля
	_, status = authenticateUser(username, "wrongPass")
	assert.Equal(t, http.StatusUnauthorized, status, "Wrong password should return 401")
}
