package main

import (
	"fmt"
	"log"
	"os"
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

// DB_URL="postgresql://postgres:password@db:5432/pc"
var db_url = os.Getenv("DB_URL")
var db *sql.DB

// Add Seed data?
// Return authenticated user
// New Image
// Delete Image
// Image - gen-short-url
// Image - List by user

func getUserIdFromToken(db *sql.DB, token string) (string, error) {
	fmt.Printf("Attempting to get user-id for token %s\n", token)
	query := fmt.Sprintf(`SELECT id
FROM users
WHERE token = '%s'
LIMIT 1`,
		token)
	// TODO: Check the token has not expired.

	rows, err := db.Query(query)
	if err != nil {
		log.Print(err)
		return "", err
	}
	defer rows.Close()

	var userid string
	for rows.Next() {
		if err := rows.Scan(&userid); err != nil {
			log.Print(err)
			return "", err
		}
	}

	return userid, nil
}

func saveNewImage(db *sql.DB, userID, href string) (err error) {
	fmt.Printf("Inserting %s for %s\n", href, userID)
	sql_query := fmt.Sprintf(`INSERT INTO images (href, user_id, created_at, updated_at)
VALUES ('%s', '%s', now(), now())`,
		href,
		userID)

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
	userId, err := getUserIdFromToken(db, token)
	if err != nil {
		c.JSON(401, gin.H{
			"error": "Invalid token",
		})
		return
	}

	// Get href string from req-body
	var reqBody map[string]interface{}
	if err := c.BindJSON(&reqBody); err != nil {
		c.JSON(400, gin.H{
			"error": "Can't parse request body",
		})
		return
	}
	href, ok := reqBody["href"].(string)
	if ok == false {
		c.JSON(400, gin.H{
			"error": "Can't parse 'href' from request body",
		})
		return
	}

	// Save into DB
	err = saveNewImage(db, userId, href)
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
	var err error
	db, err = sql.Open("postgres", db_url)
	if err != nil {
		log.Fatal(err)
	}

	r := gin.Default()
	r.GET("/ping", pong)
	r.POST("/images", createNewImage)
	r.Run() // listen and serve on 0.0.0.0:8080
}
