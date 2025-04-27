package Trakt

import (
	"fmt"
	"io"
	"net/http"
)

func SearchMovie() {
	client := &http.Client{}

	fmt.Print("Enter movie query: ")
	var query string
	fmt.Scanln(&query)

	URL := fmt.Sprintf("https://api.trakt.tv/search/movie?query=%s&extended=images", query)

	req, _ := http.NewRequest("GET", URL, nil)

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("trakt-api-version", "2")
	req.Header.Add("trakt-api-key", "9ddb03b27e491a11ac4447bd5f37b1246e78f29e3df1e2ccd8ae1d802cbf466e")

	resp, err := client.Do(req)

	if err != nil {
		fmt.Println("Errored when sending request to the server")
		return
	}

	defer resp.Body.Close()
	resp_body, _ := io.ReadAll(resp.Body)

	fmt.Println(resp.Status)
	fmt.Println(string(resp_body))
}
