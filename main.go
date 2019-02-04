package main

import (
	"fmt"
	"log"
	// "os"
	// "strconv"
	// "errors"
	"time"
	"database/sql"

	_ "github.com/lib/pq"
	uuid "github.com/pborman/uuid"
	"github.com/gin-gonic/gin"
)

type Image struct {
	UserId uuid.UUID
	Href string
	Shortlink string
}

type User struct {
	Email string
	Username string
	Token string
	TokenExpiry time.Time
}

var db *sql.DB

// Add Seed data?
// Return authenticated user
// New Image
// Delete Image
// Image - gen-short-url
// Image - List by user

func saveNewImage(db *sql.DB, userID, href string) (err error) {
	fmt.Printf("Inserting %s for %s\n", href, userID)
	sql_query := fmt.Sprintf(`INSERT INTO images (href, user_id)
VALUES ('href', 'userID')`)

	rows, err := db.Query(sql_query)
	if err != nil {
		log.Print(err)
		return err
	}
	rows.Close()
	return nil
}

func createNewImage(c *gin.Context) {
	// Get user-id from token
	// Get href string from req-body

	// Save into DB
	err := saveNewImage
	if err != nil {
		c.JSON(422, gin.H{
			"error": "Couldn't save image to DB",
		})
	}

	c.JSON(201, gin.H{})
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
