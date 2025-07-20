package main

import (
	"awesomeProject2/config"
	"awesomeProject2/internal/auth"
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"os"
)

func main() {
	config.LoadEnv()

	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		os.Getenv("DB_HOST"), os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"), os.Getenv("DB_NAME"), os.Getenv("DB_PORT"),
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("Failed to connect to DB: " + err.Error())
	}

	db.AutoMigrate(&auth.User{})

	r := gin.Default()

	repo := &auth.PgUserRepo{DB: db}
	authService := auth.NewAuthService(repo)
	auth.RegisterRoutes(r, authService)

	r.Run(":8080")
}
