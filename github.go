package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"
)

const githubURL = "https://api.github.com/app/installations"

//  Get installation token from github.
//
//  curl -i -X POST -H "Authorization: Bearer $TOKEN" -H "Accept: " https://api.github.com/app/installations/19557896/access_tokens
// {
//   "token": "asdf",
//   "expires_at": "2021-09-17T14:00:44Z",
//   "permissions": {
//     "contents": "read",
//     "metadata": "read",
//     "pull_requests": "write"
//   },
//   "repository_selection": "selected"
// }
func getInstallationToken(appInstID, token string) (string, error) {
	u := strings.Join([]string{githubURL, appInstID, "access_tokens"}, "/")

	var resBody struct {
		Token       string    `json:"token"`
		ExpiresAt   time.Time `json:"expires_at"`
		Permissions struct {
			Contents     string `json:"contents"`
			Metadata     string `json:"metadata"`
			PullRequests string `json:"pull_requests"`
		} `json:"permissions"`
		RepositorySelection string `json:"repository_selection"`
	}

	client := &http.Client{}

	req, err := http.NewRequest("POST", u, nil)
	if err != nil {
		return "", err
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	req.Header.Set("Accept", "application/vnd.github.v3+json")

	res, err := client.Do(req)
	if err != nil {
		return "", err
	}

	b, _ := ioutil.ReadAll(res.Body)

	if res.StatusCode < 200 || res.StatusCode > 300 {
		fmt.Println(res.StatusCode)
		log.Println(string(b))
		return "", fmt.Errorf("error response. status: %s, msg: %s", res.Status, string(b))
	}

	if err := json.Unmarshal(b, &resBody); err != nil {
		return "", err
	}

	return resBody.Token, nil
}
