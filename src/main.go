package main

import (
	"log"
	"os"
	"os/signal"
	"strconv"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/jerskisnow/Artemis-Bot/shards"
	"github.com/jerskisnow/Artemis-Bot/src/handlers"
	"github.com/jerskisnow/Artemis-Bot/src/utils"
	"github.com/joho/godotenv"
)

var (
	Mgr *shards.Manager
)

func main() {
	envEx := godotenv.Load(".env")
	if envEx != nil {
		log.Fatalf("[ERROR] Failed loading environment file: %s", envEx)
	}

	// ==========================================
	ce, ex := strconv.Atoi(os.Getenv("CACHE_EXPIRY"))
	if ex != nil {
		log.Fatalf("[ERROR] Could not parse cache expiry time: %s", ex)
	}

	utils.SetupCache(time.Duration(ce) * time.Minute)
	utils.SetupFirebase("firebase-credentials.json")

	// ==========================================
	Mgr, botEx := shards.New("Bot " + os.Getenv("BOT_TOKEN"))
	if botEx != nil {
		log.Fatalf("[ERROR] Session creation failed: %s", botEx)
	}

	Mgr.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
		// TODO: Change this to: Watching the wind guide my arrows
		s.UpdateGameStatus(0, "with bow and arrow")
	})

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

	log.Println("[INFO] Starting sharding manager...")
	shardEx := Mgr.Start()
	if shardEx != nil {
		log.Fatalf("[ERROR] Couldn't start the sharding manager.")
	}

	// ==========================================
	if os.Getenv("PRODUCTION") == "0" {
		handlers.RegisterCommands(Mgr, os.Getenv("GUILD_ID"))
	} else {
		handlers.RegisterCommands(Mgr, "")
	}

	// ==========================================
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	<-stop

	log.Println("[INFO] Stopping sharding manager...")
	Mgr.Shutdown()
	log.Println("[SUCCESS] Safe-Shutdown completed.")
}
