package auth

import (
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"testing"
)

func setupTestAuthService() *AuthService {
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		panic("in-memory DB başlatılamadı: " + err.Error())
	}
	_ = db.AutoMigrate(&User{})
	repo := &PgUserRepo{DB: db}
	return NewAuthService(repo)
}

func TestRegisterAndLogin_Success(t *testing.T) {
	service := setupTestAuthService()

	// Register
	err := service.Register("test@example.com", "123456")
	assert.NoError(t, err)

	// Login
	user, err := service.Login("test@example.com", "123456")
	assert.NoError(t, err)
	assert.Equal(t, "test@example.com", user.Email)
}

func TestRegister_ExistingUser(t *testing.T) {
	service := setupTestAuthService()

	_ = service.Register("test@example.com", "123456")

	// Try again
	err := service.Register("test@example.com", "123456")
	assert.ErrorIs(t, err, ErrUserExists)
}

func TestLogin_InvalidPassword(t *testing.T) {
	service := setupTestAuthService()

	_ = service.Register("test@example.com", "123456")

	_, err := service.Login("test@example.com", "wrongpass")
	assert.ErrorIs(t, err, ErrInvalidCredentials)
}

func TestLogin_UserNotFound(t *testing.T) {
	service := setupTestAuthService()

	_, err := service.Login("notfound@example.com", "123456")
	assert.ErrorIs(t, err, ErrInvalidCredentials)
}
