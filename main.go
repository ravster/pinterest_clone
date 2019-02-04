package main

import (
	// "fmt"
	// "os"
	// "strconv"

	"time"

//	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/gin-gonic/gin"
)

type Image struct {
	Href string
	Shortlink string
}

type User struct {
	Email string
	Username string
	Token string
	TokenExpiry time
}

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
