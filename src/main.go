package main

import (
	"log"
	"os"
	"os/signal"

	"github.com/bwmarrin/discordgo"
	"github.com/jerskisnow/Suggestions/src/handlers"
	"github.com/jerskisnow/Suggestions/src/utils"
	"github.com/joho/godotenv"
	"github.com/servusdei2018/shards"
)

var (
	Mgr *shards.Manager
)

func main() {
	envEx := godotenv.Load(".env")
	if envEx != nil {
		log.Fatalf("[ERROR] Couldn't load the environment file.")
		return
	}

	Mgr, botEx := shards.New("Bot " + os.Getenv("BOT_TOKEN"))
	if botEx != nil {
		log.Fatalf("[ERROR] Couldn't create a session.")
		return
	}

	// Add the event handlers
	Mgr.AddHandler(func(s *discordgo.Session, e *discordgo.Connect) {
		log.Printf("[INFO] Shard #%v connected.\n", s.ShardID)
	})

	Mgr.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if i.Type == 2 { // ApplicationCommand
			handlers.LinkCommand(s, i)
		} else if i.Type == 3 { // MessageComponent
			handlers.LinkButton(s, i)
		}
	})

	// Setup redis
	utils.SetupCache()

	log.Println("[INFO] Starting sharding manager...")
	shardEx := Mgr.Start()
	if shardEx != nil {
		log.Fatalf("[ERROR] Couldn't start the sharding manager.")
		return
	}

	// Create commands
	if os.Getenv("PRODUCTION") == "0" {
		handlers.RegisterCommands(Mgr, os.Getenv("GUILD_ID"))
	} else {
		handlers.RegisterCommands(Mgr, "")
	}

	// On shutdown handles the stuff below
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	<-stop

	log.Println("[INFO] Stopping sharding manager...")
	Mgr.Shutdown()
	log.Println("[SUCCESS] Safe-Shutdown completed.")
}
