package db

import (
	"os"
	"database/sql"
	_ "github.com/lib/pq"
	"log"
	"fmt"
	"errors"
)

// DB_URL="postgresql://postgres:password@db:5432/pc"
var db_url = os.Getenv("DB_URL")
var DB *sql.DB

func Connect() {
	if DB != nil {
		DB.Close()
	}

	var err error
	log.Println("connecting to postgres database... ")
	DB, err = sql.Open("postgres", db_url)

	if err != nil {
		log.Panicf(`failed to connect database. err="%v"`, err)
	}
}

func GetUserIdFromToken(token string) (string, error) {
	fmt.Printf("Attempting to get user-id for token %s\n", token)
	query := fmt.Sprintf(`SELECT id
FROM users
WHERE token = '%s'
LIMIT 1`,
		token)
	// TODO: Check the token has not expired.

	rows, err := DB.Query(query)
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
	if userid == "" {
		return "", errors.New("Token not found")
	}

	return userid, nil
}

func SaveNewImage(userID, href string) (err error) {
	fmt.Printf("Inserting %s for %s\n", href, userID)
	sql_query := fmt.Sprintf(`INSERT INTO images (href, user_id, created_at, updated_at)
VALUES ('%s', '%s', now(), now())`,
		href,
		userID)

	rows, err := DB.Query(sql_query)
	if err != nil {
		log.Print(err)
		return err
	}
	rows.Close()
	return nil
}

func MarkImageDeleted(userId, imageId string) error {
	fmt.Printf("Deleting image %s\n", imageId)
	sql_query := fmt.Sprintf(`UPDATE images set deleted_at = now()
WHERE user_id = '%s' AND id = '%s'`,
		userId,
		imageId,
	)

	rows, err := DB.Query(sql_query)
	if err != nil {
		log.Print(err)
		return err
	}
	rows.Close()
	return nil
}

func ListImagesForUser(userId string) ([]string, error) {
	var hrefs []string

	fmt.Printf("Getting hrefs of all undeleted images for user-id %s\n", userId)
	query := fmt.Sprintf(`SELECT href
FROM images
WHERE user_id = '%s' AND deleted_at IS NULL
ORDER BY created_at DESC`,
		userId)

	rows, err := DB.Query(query)
	if err != nil {
		log.Print(err)
		return hrefs, err
	}
	defer rows.Close()

	var href string
	for rows.Next() {
		if err := rows.Scan(&href); err != nil {
			log.Print(err)
			return hrefs, err
		}

		hrefs = append(hrefs, href)
	}

	return hrefs, nil
}
