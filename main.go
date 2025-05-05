package main

import (
	"fmt"

	// "github.com/ChrolloKyber/Traktor/Bot"

	"github.com/ChrolloKyber/Traktor/Bot"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading environment variables: ", err)
	}
	// Trakt.AccessToken()
	Bot.Run()
	/* result := Trakt.SearchContent("movie", "Tron")
	fmt.Println(result.Movie.Title)
	fmt.Println(result.Movie.Year)

	values := reflect.ValueOf(result.Movie.IDs)
	for k := range values.NumField() {
		fmt.Printf("%s: %v\n", values.Type().Field(k).Name, values.Field(k).Interface())
	} */
}
