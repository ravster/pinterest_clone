package github

import (
	"net/http"
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
)

func GetAccessTokenFromGithubLogin(code string) (bool, string) {
	client := &http.Client{}
	clientId := "0028f2b81b2b5aa770b3"
	clientSecret := "1db7aee5c6488d7a0b8261fb7ecca95537c8d6cb"

	bodyString := fmt.Sprintf("code=%s&client_id=%s&client_secret=%s", code, clientId, clientSecret)
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
		return false, "Couldn't parse JSON from Github"
	}
	if respFromGH["access_token"] == "" {
		return false, "Couldn't get an access-token from GitHub"
	}

	return true, ""
}
