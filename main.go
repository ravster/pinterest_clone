package main

import (
	"github.com/gin-gonic/gin"

	db "github.com/ravster/pinterest_clone/db"
	github "github.com/ravster/pinterest_clone/github"
)

type userGetter func (string) (string, error)

func getUserIdFromToken(userGetterFunc userGetter, token string) (string, string) {
	if token == "" {
		return "", "Missing Authorization token"
	}

	userId, err := userGetterFunc(token)
	if err != nil {
		return "", "Invalid token"
	}

	return userId, ""
}

func createNewImage(c *gin.Context) {
	token := c.GetHeader("Authorization")
	userId, errstring := getUserIdFromToken(db.GetUserIdFromToken, token)
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
	userId, errstring := getUserIdFromToken(db.GetUserIdFromToken, token)
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

func getMyImages(c *gin.Context) {
	token := c.GetHeader("Authorization")
	userId, errstring := getUserIdFromToken(db.GetUserIdFromToken, token)
	if errstring != "" {
		c.JSON(401, gin.H{
			"error": errstring,
		})
		return
	}

	hrefs, err := db.ListImagesForUser(userId)
	if err != nil {
		c.JSON(500, gin.H{
			"error": "Couldn't get the list of images",
		})
		return
	}

	c.JSON(200, hrefs)
}

func getUserImages(c *gin.Context) {
	userId := c.Param("userId")

	hrefs, err := db.ListImagesForUser(userId)
	if err != nil {
		c.JSON(500, gin.H{
			"error": "Couldn't get the list of images",
		})
		return
	}

	c.JSON(200, hrefs)
}

type GHResponse struct {
	AccessToken string `json:access_token`
}

func loginFromGitHub(c *gin.Context) {
	accessCode, ok := c.GetQuery("code")
	if ok == false {
		c.JSON(422, gin.H{
			"error": "No code found from GitHub callback",
		})
		return
	}

	ok, errString := github.GetAccessTokenFromGithubLogin(accessCode)
	if errString != "" {
		c.JSON(500, gin.H{
			"error": errString,
		})
		return
	}

	userId := c.Param("userId")

	token, err := db.CreateToken(userId)
	if err != nil {
		c.JSON(500, gin.H{
			"error": "Couldn't create token for the application",
		})
		return
	}

	c.JSON(201, gin.H{
		"access_token": token,
	})
}

func main() {
	db.Connect()

	r := gin.Default()

	r.POST("/images", createNewImage)
	r.DELETE("/images/:id", deleteImage)
	r.GET("/images", getMyImages)

	r.GET("/images/:userId", getUserImages)
	r.GET("/success_GH_authn_callback/:userId", loginFromGitHub)

	r.Run() // listen and serve on 0.0.0.0:8080
}
