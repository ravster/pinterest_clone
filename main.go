package main

import (
	// "fmt"
	// "os"
	// "strconv"

	"time"

	//	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/pborman/uuid"
	"github.com/gin-gonic/gin"
)

type Image struct {
	UserId UUID
	Href string
	Shortlink string
}

type User struct {
	Email string
	Username string
	Token string
	TokenExpiry time
}

// Add Seed data?
// Return authenticated user
// New Image
// Delete Image
// Image - gen-short-url
// Image - List by user

func pong(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "pong",
	})
}

func main() {
	r := gin.Default()
	r.GET("/ping", pong)
	r.Run() // listen and serve on 0.0.0.0:8080
}
