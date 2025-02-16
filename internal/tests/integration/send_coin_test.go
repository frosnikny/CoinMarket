package integration

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func createUser(username, password string) string {
	authPayload := map[string]string{
		"username": username,
		"password": password,
	}
	body, _ := json.Marshal(authPayload)

	req, _ := http.NewRequest("POST", "/api/auth", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	testRouter.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		log.Fatalf("Ошибка создания пользователя %s: %s", username, w.Body.String())
	}

	var authResp map[string]string
	json.Unmarshal(w.Body.Bytes(), &authResp)
	return authResp["token"]
}
func TestSendCoins(t *testing.T) {
	// 1. Создаем пользователей и получаем JWT-токены
	senderToken := createUser("senderUser", "password123")
	receiverToken := createUser("receiverUser", "password123")

	// 2. Получаем баланс отправителя и получателя перед переводом
	req, _ := http.NewRequest("GET", "/api/info", nil)
	req.Header.Set("Authorization", "Bearer "+senderToken)
	w := httptest.NewRecorder()
	testRouter.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	var senderInfo map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &senderInfo)
	senderCoinsBefore := int(senderInfo["coins"].(float64))

	req, _ = http.NewRequest("GET", "/api/info", nil)
	req.Header.Set("Authorization", "Bearer "+receiverToken)
	w = httptest.NewRecorder()
	testRouter.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	var receiverInfo map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &receiverInfo)
	receiverCoinsBefore := int(receiverInfo["coins"].(float64))

	// 3. Отправляем 50 монет от senderUser к receiverUser
	sendPayload := map[string]interface{}{
		"toUser": "receiverUser",
		"amount": 50,
	}
	body, _ := json.Marshal(sendPayload)

	req, _ = http.NewRequest("POST", "/api/sendCoin", bytes.NewBuffer(body))
	req.Header.Set("Authorization", "Bearer "+senderToken)
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	testRouter.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	// 4. Проверяем баланс после перевода
	req, _ = http.NewRequest("GET", "/api/info", nil)
	req.Header.Set("Authorization", "Bearer "+senderToken)
	w = httptest.NewRecorder()
	testRouter.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	json.Unmarshal(w.Body.Bytes(), &senderInfo)
	senderCoinsAfter := int(senderInfo["coins"].(float64))
	assert.Equal(t, senderCoinsBefore-50, senderCoinsAfter, "Sender should have 50 coins less")

	req, _ = http.NewRequest("GET", "/api/info", nil)
	req.Header.Set("Authorization", "Bearer "+receiverToken)
	w = httptest.NewRecorder()
	testRouter.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	json.Unmarshal(w.Body.Bytes(), &receiverInfo)
	receiverCoinsAfter := int(receiverInfo["coins"].(float64))
	assert.Equal(t, receiverCoinsBefore+50, receiverCoinsAfter, "Receiver should have 50 coins more")

	// 5. Проверяем, что транзакция записана в истории
	sentHistory := senderInfo["coinHistory"].(map[string]interface{})["sent"].([]interface{})
	receivedHistory := receiverInfo["coinHistory"].(map[string]interface{})["received"].([]interface{})

	foundSent := false
	log.Println(senderInfo)
	log.Println(sentHistory)
	for _, txn := range sentHistory {
		txnMap := txn.(map[string]interface{})
		if txnMap["ToUser"] == "receiverUser" && int(txnMap["Amount"].(float64)) == 50 {
			foundSent = true
			break
		}
	}
	assert.True(t, foundSent, "Transaction should be in sender's history")

	foundReceived := false
	log.Println(receivedHistory)
	for _, txn := range receivedHistory {
		txnMap := txn.(map[string]interface{})
		if txnMap["FromUser"] == "senderUser" && int(txnMap["Amount"].(float64)) == 50 {
			foundReceived = true
			break
		}
	}
	assert.True(t, foundReceived, "Transaction should be in receiver's history")
}
