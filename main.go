package main

import (
	"encoding/json"
	"net/http"
	"bytes"
	"strings"
	"fmt"

	"github.com/gin-gonic/gin"

	db "github.com/ravster/pinterest_clone/db"
)

// Add Seed data?
// Return authenticated user
// Image - gen-short-url

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
	token := c.GetHeader("Authorization")
	userId, errstring := getUserIdFromToken(token)
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

	client := &http.Client{}
	clientId := "0028f2b81b2b5aa770b3"
	clientSecret := "1db7aee5c6488d7a0b8261fb7ecca95537c8d6cb"
	bodyString := fmt.Sprintf("code=%s&client_id=%s&client_secret=%s", accessCode, clientId, clientSecret)
	reqBody := strings.NewReader(bodyString)
	req, _ := http.NewRequest("POST", "https://github.com/login/oauth/access_token", reqBody)
	req.Header.Add("Accept", "application/json")

	resp, _ := client.Do(req)

	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	newStr := buf.String()

	var respFromGH map[string]string

	err := json.Unmarshal([]byte(newStr), &respFromGH)
	if err != nil {
		c.JSON(500, gin.H{
			"error": "Couldn't parse JSON from Github",
		})
		return
	}
	if respFromGH["access_token"] == "" {
		c.JSON(400, gin.H{
			"error": "Couldn't get an access-token from GitHub",
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

	// At this point, we've confirmed that the github authentication passed.  Now, how do we ensure that we log in as the right user?
	// Just set the user's UUID in a subdir of the path of the callback URL registered with GH.
	// https://github.com/login/oauth/authorize?scope=user:email&client_id=0028f2b81b2b5aa770b3&redirect_uri=http://localhost:8080/success_GH_authn_callback/0208be54-e388-4ed1-b435-c2b063cce9c1
}

func main() {
	db.Connect()

	r := gin.Default()

	// curl -H "Authorization: foo" -XPOST -d '{"href": "http://foo.com"}' localhost:8080/images
	r.POST("/images", createNewImage)
	// curl -H "Authorization: foo" -XDELETE localhost:8080/images/621179b9-a872-4452-aa01-415507ff9b44
	r.DELETE("/images/:id", deleteImage)
	r.GET("/images", getMyImages) // curl -H "Authorization: foo" -XGET localhost:8080/images

	r.GET("/ping", pong)
	r.GET("/images/:userId", getUserImages) // curl -XGET localhost:8080/images/0208be54-e388-4ed1-b435-c2b063cce9c1
	r.GET("/success_GH_authn_callback/:userId", loginFromGitHub)

	// http://localhost:8080/success_GH_authn_callback?code=708cd6929b3c7a168f7b

	r.Run() // listen and serve on 0.0.0.0:8080
}
