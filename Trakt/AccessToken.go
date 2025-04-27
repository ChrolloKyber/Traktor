package Trakt

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
)

type AccessTokens struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

func AccessToken() AccessTokens {
	err := godotenv.Load()
	client := &http.Client{}

	clientID := os.Getenv("CLIENT_ID")
	clientSecret := os.Getenv("CLIENT_SECRET")

	bodyStruct := struct {
		Code         string `json:"code"`
		ClientID     string `json:"client_id"`
		ClientSecret string `json:"client_secret"`
	}{
		Code:         DeviceCode().DeviceCode,
		ClientID:     clientID,
		ClientSecret: clientSecret,
	}

	bodyBytes, err := json.Marshal(bodyStruct)
	if err != nil {
		fmt.Println("Error marshalling JSON")
		fmt.Println(err)
	}

	body := []byte(bodyBytes)

	var req *http.Request
	var resp *http.Response

	i := 0
	for {
		if i == 60 {
			break
		}

		req, _ = http.NewRequest("POST", "https://api.trakt.tv/oauth/device/token", bytes.NewBuffer(body))

		req.Header.Add("Content-Type", "application/json")

		resp, err = client.Do(req)
		if err != nil {
			fmt.Println("Error making request")
			fmt.Println(err)
			return AccessTokens{}
		}
		defer resp.Body.Close()

		fmt.Println(resp.StatusCode)
		if resp.StatusCode == 200 || resp.StatusCode == 409 || resp.StatusCode == 410 || resp.StatusCode == 418 {
			break
		}
		resp_body, _ := io.ReadAll(resp.Body)
		fmt.Println(string(resp_body))
		time.Sleep(time.Second * time.Duration(10))
		i++
	}

	if err != nil {
		fmt.Println("Errored when sending request to the server")
	}

	defer resp.Body.Close()
	resp_body, _ := io.ReadAll(resp.Body)

	var accessToken AccessTokens
	err = json.Unmarshal(resp_body, &accessToken)

	if err != nil {
		fmt.Println("Error unmarshalling the JSON values")
		fmt.Println(err)
	}

	return accessToken
}
