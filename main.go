package main

import (
	"fmt"

	"github.com/ChrolloKryber/Traktor/Bot"
	// "github.com/ChrolloKryber/Traktor/Trakt"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading environment variables: ", err)
	}
	// Trakt.AccessToken()
	Bot.Run()
}
