package Bot

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"

	"github.com/ChrolloKyber/Traktor/Trakt"
	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
)

var (
	GuildID        = flag.String("guild", "1365970657665482782", "Test guild ID. If not passed - bot registers commands globally")
	RemoveCommands = flag.Bool("rmcmd", true, "Remove all commands after shutdowning or not")
)

var s *discordgo.Session

func checkErr(errorMessage string, err error) {
	if err != nil {
		fmt.Printf("%s: %v", errorMessage, err)
	}
}

func init() {
	err := godotenv.Load()
	checkErr("Error loading environment variables", err)

	s, err = discordgo.New("Bot " + os.Getenv("DISCORD_TOKEN"))
	checkErr("Error logging in", err)
}

var (
	defaultMemberPermissions int64 = discordgo.PermissionManageServer

	commands = []*discordgo.ApplicationCommand{
		{
			Name:        "search",
			Description: "Search for movies or series",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "type",
					Required:    true,
					Description: "Choose between movies and series",
					Choices: []*discordgo.ApplicationCommandOptionChoice{
						{
							Name:  "Movies",
							Value: "movie",
						},
						{
							Name:  "Show",
							Value: "show",
						},
						{
							Name:  "Person",
							Value: "person",
						},
						{
							Name:  "Episode",
							Value: "episode",
						},
						{
							Name:  "List",
							Value: "list",
						},
					},
				},
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "query",
					Description: "Enter your query",
					Required:    true,
				},
			},
		},
	}

	commandHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
		"search": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			opts := i.ApplicationCommandData().Options

			contentType := opts[0].StringValue()
			query := opts[1].StringValue()

			data, err := Trakt.SearchContent(contentType, query)

			if err != nil {
				fmt.Println("Error searching content: ", err)
			}

			var response string
			if contentType == "movie" {
				response = fmt.Sprintf("%s: %s\nYear: %d\nIMDB: https://imdb.com/title/%s\n", strings.Title(contentType), data.Movie.Title, data.Movie.Year, data.Movie.IDs.IMDB)
			} else if contentType == "show" {
				response = fmt.Sprintf("%s: %s\nYear: %d\nIMDB: https://imdb.com/title/%s\n", strings.Title(contentType), data.Show.Title, data.Show.Year, data.Show.IDs.IMDB)
			} else if contentType == "person" {
				response = fmt.Sprintf("Name: %s\nIMDB: https://trakt.tv/people/%s\n", data.Person.Name, data.Person.IDs.Slug)
			} else if contentType == "list" {
				response = fmt.Sprintf("List: %s\nCreated by: [%s](https://trakt.tv/users/%s)\nTrakt Link: %s\n", data.List.Name, data.List.User.Name, data.List.User.IDs.Slug, data.List.ShareLink)
			} else if contentType == "episode" {
				response = fmt.Sprintf("Episode title: %s\nShow: %s (Season: %d, Episode: %d)\nLink: https://trakt.tv/shows/%s/seasons/%d/episodes/%d\n", data.Episode.Title, data.Show.Title, data.Episode.Season, data.Episode.Number, data.Show.IDs.Slug, data.Episode.Season, data.Episode.Number)
			}

			if err != nil {
				response = fmt.Sprintf("Error returned by the search function: %v", err)
			}

			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: response,
				},
			})
		},
	}
)

func init() {
	s.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if h, ok := commandHandlers[i.ApplicationCommandData().Name]; ok {
			h(s, i)
		}
	})
}

func Run() {
	s.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
		log.Printf("Logged in as: %v#%v", s.State.User.Username, s.State.User.Discriminator)
	})
	err := s.Open()
	if err != nil {
		log.Fatalf("Cannot open the session: %v", err)
	}

	log.Println("Adding commands...")
	registeredCommands := make([]*discordgo.ApplicationCommand, len(commands))
	for i, v := range commands {
		cmd, err := s.ApplicationCommandCreate(s.State.User.ID, *GuildID, v)
		if err != nil {
			log.Panicf("Cannot create '%v' command: %v", v.Name, err)
		}
		registeredCommands[i] = cmd
	}

	defer s.Close()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	log.Println("Press Ctrl+C to exit")
	<-stop

	if *RemoveCommands {
		log.Println("Removing all commands...")
		// Fetch all commands to make sure we get everything
		// registeredCommands, err := s.ApplicationCommands(s.State.User.ID, *GuildID)
		// if err != nil {
		// log.Fatalf("Could not fetch commands: %v", err)
		// }

		log.Printf("Found %d commands to delete", len(registeredCommands))

		for _, v := range registeredCommands {
			log.Printf("Deleting command: %s (ID: %s)", v.Name, v.ID)
			err := s.ApplicationCommandDelete(s.State.User.ID, *GuildID, v.ID)
			if err != nil {
				log.Printf("Error deleting command '%v': %v", v.Name, err)
			} else {
				log.Printf("Successfully deleted command: %s", v.Name)
			}
		}
	}

	log.Println("Gracefully shutting down.")
}
