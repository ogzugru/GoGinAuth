package auth

import (
	"awesomeProject2/internal/pkg/utils"
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

const hashedPassword = "$2a$14$zF66Lgn6lBIDRZxCrFhhwe3fYGoAK0DFeXNdvTqMoU2ykTSttP08i"

func setupTestRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.Default()

	hashed, err := utils.HashPassword("123456")
	if err != nil {
		panic("şifre hashlenemedi: " + err.Error())
	}

	fakeUserStore["test@example.com"] = User{
		ID:       99,
		Email:    "test@example.com",
		Password: hashed,
	}

	RegisterRoutes(r)
	return r
}

func TestLogin_Success(t *testing.T) {
	router := setupTestRouter()

	body := map[string]string{
		"email":    "test@example.com",
		"password": "123456",
	}
	jsonBody, _ := json.Marshal(body)

	req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	assert.Contains(t, response, "token")
	assert.Contains(t, response, "user")
}

func TestLogin_InvalidPassword(t *testing.T) {
	router := setupTestRouter()

	body := map[string]string{
		"email":    "test@example.com",
		"password": "wrongpassword",
	}
	jsonBody, _ := json.Marshal(body)

	req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestLogin_UserNotFound(t *testing.T) {
	router := setupTestRouter()

	body := map[string]string{
		"email":    "notfound@example.com",
		"password": "123456",
	}
	jsonBody, _ := json.Marshal(body)

	req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}
