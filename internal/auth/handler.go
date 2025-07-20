package auth

import (
	"awesomeProject2/internal/middleware"
	"awesomeProject2/internal/pkg/utils"
	"github.com/gin-gonic/gin"
	"net/http"
)

var fakeUserStore = map[string]User{}
var currentID int64 = 1

func RegisterRoutes(r *gin.Engine) {
	r.POST("/register", register)
	r.POST("/login", login)

	protected := r.Group("/")
	protected.Use(middleware.AuthRequired())
	protected.GET("/me", me)
}

func register(c *gin.Context) {
	var input RegisterInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if _, exists := fakeUserStore[input.Email]; exists {
		c.JSON(http.StatusConflict, gin.H{"error": "User already exists"})
		return
	}

	hashed, err := utils.HashPassword(input.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	user := User{
		ID:       currentID,
		Email:    input.Email,
		Password: hashed,
	}
	currentID++
	fakeUserStore[user.Email] = user

	c.JSON(http.StatusOK, gin.H{"message": "User registered"})
}

func login(c *gin.Context) {
	var input LoginInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	storedUser, exists := fakeUserStore[input.Email]
	if !exists || !utils.CheckPasswordHash(input.Password, storedUser.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	token, err := utils.GenerateJWT(storedUser.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token": token,
		"user": gin.H{
			"id":    storedUser.ID,
			"email": storedUser.Email,
		},
	})
}

func me(c *gin.Context) {
	userIDRaw, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	userID := userIDRaw.(int64)

	for _, user := range fakeUserStore {
		if user.ID == userID {
			c.JSON(http.StatusOK, gin.H{
				"id":    user.ID,
				"email": user.Email,
			})
			return
		}
	}

	c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
}
