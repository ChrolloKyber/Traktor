package Trakt

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type SearchResult struct {
	Type    string   `json:"type"`
	Score   float64  `json:"score"`
	Movie   *Movie   `json:"movie,omitempty"`
	Show    *Show    `json:"show,omitempty"`
	Person  *Person  `json:"person,omitempty"`
	List    *List    `json:"list,omitempty"`
	Episode *Episode `json:"episode,omitempty"`
}

type Movie struct {
	Title string `json:"title"`
	Year  int    `json:"year"`
	IDs   IDs    `json:"ids"`
}

type Show struct {
	Title string `json:"title"`
	Year  int    `json:"year"`
	IDs   IDs    `json:"ids"`
}

type Episode struct {
	Season int    `json:"season"`
	Number int    `json:"number"`
	Title  string `json:"title"`
	IDs    IDs    `json:"ids"`
	Show   *Show  `json:"show,omitempty"`
}

type Person struct {
	Name string `json:"name"`
	IDs  IDs    `json:"ids"`
}

type IDs struct {
	Trakt int    `json:"trakt,omitempty"`
	Slug  string `json:"slug,omitempty"`
	IMDB  string `json:"imdb,omitempty"`
	TMDB  int    `json:"tmdb,omitempty"`
	TVDB  int    `json:"tvdb,omitempty"`
}

type User struct {
	Username string `json:"username"`
	Private  bool   `json:"private"`
	Name     string `json:"name"`
	VIP      bool   `json:"vip"`
	VIPEP    bool   `json:"vip_ep"`
	IDs      struct {
		Slug string `json:"slug"`
	} `json:"ids"`
}

type List struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Privacy     string `json:"privacy"`
	ShareLink   string `json:"share_link"`
	Type        string `json:"type"`
	ItemCount   int    `json:"item_count"`
	Likes       int    `json:"likes"`
	IDs         struct {
		Trakt int    `json:"trakt"`
		Slug  string `json:"slug"`
	} `json:"ids"`
	User User `json:"user"`
}

func SearchContent(contentType, query string) SearchResult {
	client := &http.Client{}

	URL := fmt.Sprintf("https://api.trakt.tv/search/%s?query=%s&extended=images", contentType, query)

	req, _ := http.NewRequest("GET", URL, nil)

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("trakt-api-version", "2")
	req.Header.Add("trakt-api-key", "9ddb03b27e491a11ac4447bd5f37b1246e78f29e3df1e2ccd8ae1d802cbf466e")

	resp, err := client.Do(req)

	if err != nil {
		fmt.Println("Errored when sending request to the server")
		return SearchResult{}
	}

	var results []SearchResult
	defer resp.Body.Close()
	resp_body, _ := io.ReadAll(resp.Body)

	err = json.Unmarshal(resp_body, &results)

	fmt.Println(resp.Status)
	fmt.Println(string(resp_body))
	return results[0]
}
