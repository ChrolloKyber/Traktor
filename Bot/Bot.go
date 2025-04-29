package Bot

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"

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
							Name:  "Series",
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

			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: fmt.Sprintf("Content Type: %s\nQuery: %s\n", contentType, query),
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
