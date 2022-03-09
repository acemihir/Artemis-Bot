package main

import (
	"log"
	"os"
	"os/signal"

	"github.com/bwmarrin/discordgo"
	"github.com/jerskisnow/Suggestions/src/handlers"
	"github.com/joho/godotenv"
)

var s *discordgo.Session

func main() {
	envEx := godotenv.Load(".env")
	if envEx != nil {
		log.Fatalf("Couldn't load the environment file.")
	}

	s, botEx := discordgo.New("Bot " + os.Getenv("BOT_TOKEN"))
	if botEx != nil {
		log.Fatalf("Couldn't create a session.")
	}

	// Add the event handlers
	s.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
		log.Printf("Logged in as: %v#%v", s.State.User.Username, s.State.User.Discriminator)
	})

	s.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if i.Type == 2 { // ApplicationCommand
			handlers.LinkCommand(s, i)
		} else if i.Type == 3 { // MessageComponent
			handlers.LinkButton(s, i)
		}
	})

	// Open the actual session
	sesEx := s.Open()
	if sesEx != nil {
		log.Fatalf("Couldn't open a session.")
	}

	// Create commands
	if os.Getenv("PRODUCTION") == "0" {
		handlers.RegisterCommands(s, os.Getenv("GUILD_ID"))
	} else {
		handlers.RegisterCommands(s, "")
	}

	defer s.Close()

	// On shutdown handles the stuff below
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	<-stop

	// TODO: Do something...

	log.Println("Safe-Shutdown completed.")
}
