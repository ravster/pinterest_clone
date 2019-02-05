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

func getUserIdFromToken(db *sql.DB, token string) (string, err) {
	fmt.Printf("Attempting to get user-id for token %s\n", token)
	query := fmt.Sprintf(`SELECT id
FROM users
WHERE token = '%s'
LIMIT 1`,
		token)

	rows, err := db.Query(query)
	if err != nil {
		log.Print(err)
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var userid string
		if err := rows.Scan(&userid); err != nil {
			log.Print(err)
			return nil, err
		}

		return userid, nil
	}
}

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
	// Get token from context
	token := c.GetHeader("Authorization")
	if token == "" {
		c.JSON(401, gin.H{
			"error": "Missing Authorization token",
		})
		return
	}
	// Get user-id from token
	userId, err := getUserIdFromToken(token)
	if err != nil {
		c.JSON(401, gin.H{
			"error": "Invalid token",
		})
		return
	}
	// Get href string from req-body

	// Save into DB
	err := saveNewImage
	if err != nil {
		c.JSON(422, gin.H{
			"error": "Couldn't save image to DB",
		})
		return
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
