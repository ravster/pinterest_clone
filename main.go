package main

import (
	"github.com/gin-gonic/gin"

	db "github.com/ravster/pinterest_clone/db"
)

// Add Seed data?
// Return authenticated user
// Image - gen-short-url
// Image - List by user

func getUserIdFromToken(token string) (string, string) {
	if token == "" {
		return "", "Missing Authorization token"
	}

	userId, err := db.GetUserIdFromToken(token)
	if err != nil {
		return "", "Invalid token"
	}

	return userId, ""
}

func createNewImage(c *gin.Context) {
	token := c.GetHeader("Authorization")
	userId, errstring := getUserIdFromToken(token)
	if errstring != "" {
		c.JSON(401, gin.H{
			"error": errstring,
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

	if err := db.SaveNewImage(userId, href); err != nil {
		c.JSON(422, gin.H{
			"error": "Couldn't save image to DB",
		})
		return
	}

	c.JSON(201, gin.H{})
}

func deleteImage(c *gin.Context) {
	token := c.GetHeader("Authorization")
	userId, errstring := getUserIdFromToken(token)
	if errstring != "" {
		c.JSON(401, gin.H{
			"error": errstring,
		})
		return
	}

	// get id of image
	imageId := c.Param("id")

	// mark image deleted in DB
	if err := db.MarkImageDeleted(userId, imageId);  err != nil {
		c.JSON(422, gin.H{
			"error": "Couldn't delete image",
		})
		return
	}

	c.JSON(200, gin.H{})

}

func pong(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "pong",
	})
}

func getMyImages(c *gin.Context) {
	// get token
	// get userid
	// get image-hrefs
	// convert to JSON
	// HTTP response
}

func main() {
	db.Connect()

	r := gin.Default()
	r.GET("/ping", pong)
	// curl -H "Authorization: foo" -XPOST -d '{"href": "http://foo.com"}' localhost:8080/images
	r.POST("/images", createNewImage)
	// curl -H "Authorization: foo" -XDELETE localhost:8080/images/621179b9-a872-4452-aa01-415507ff9b44
	r.DELETE("/images/:id", deleteImage)
	r.GET("/images", getMyImages)
	r.Run() // listen and serve on 0.0.0.0:8080
}
