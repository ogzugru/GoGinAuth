package auth

import (
	"awesomeProject2/internal/middleware"
	"awesomeProject2/internal/pkg/utils"
	"github.com/gin-gonic/gin"
	"net/http"
)

func RegisterRoutes(r *gin.Engine, service *AuthService) {
	r.POST("/register", func(c *gin.Context) {
		var input RegisterInput
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		err := service.Register(input.Email, input.Password)
		if err != nil {
			if err == ErrUserExists {
				c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to register user"})
			}
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "User registered"})
	})

	r.POST("/login", func(c *gin.Context) {
		var input LoginInput
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		user, err := service.Login(input.Email, input.Password)
		if err != nil {
			if err == ErrInvalidCredentials {
				c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Login failed"})
			}
			return
		}

		token, err := utils.GenerateJWT(int64(user.ID))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"token": token,
			"user": gin.H{
				"id":    user.ID,
				"email": user.Email,
			},
		})
	})

	protected := r.Group("/")
	protected.Use(middleware.AuthRequired())
	protected.GET("/me", func(c *gin.Context) {
		userIDRaw, exists := c.Get("user_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}

		userID := uint(userIDRaw.(int64))
		user, err := service.Repo.FindByID(userID)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"id":    user.ID,
			"email": user.Email,
		})
	})
}
