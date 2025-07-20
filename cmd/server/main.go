package main

import (
	"awesomeProject2/internal/auth"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	auth.RegisterRoutes(r)

	r.Run(":8080")
}
