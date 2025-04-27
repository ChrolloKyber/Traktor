package Trakt

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"

	"github.com/joho/godotenv"
)

func MyMovies() {
	err := godotenv.Load()
	client := &http.Client{}

	clientID := os.Getenv("CLIENT_ID")
	accessToken := AccessToken()

	urlString := fmt.Sprintf("https://api.trakt.tv/users/amulmakkhan/watched/movies") //, start_date, days)
	bearer := fmt.Sprintf("Bearer %s", accessToken.AccessToken)

	req, _ := http.NewRequest("GET", urlString, nil)

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", bearer)
	req.Header.Add("trakt-api-version", "2")
	req.Header.Add("trakt-api-key", clientID)

	resp, err := client.Do(req)

	if err != nil {
		fmt.Println("Errored when sending request to the server")
		return
	}

	defer resp.Body.Close()
	resp_body, _ := io.ReadAll(resp.Body)

	fmt.Println(resp.Status)
	fmt.Println(string(resp_body))
	file, err := os.Create("test.json")
	if err != nil {
		fmt.Println("Error creating file")
	}
	defer file.Close()

	_, err = file.Write(resp_body)

	if err != nil {
		fmt.Println("Error writing to file: ", err)
	}
	command := exec.Command("prettier", "-w", "test.json")
	fmt.Println(command.Output())
}
