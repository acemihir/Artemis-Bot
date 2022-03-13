package main

import (
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
	utils.Cout("[INFO] Start sequence initiated.\n", utils.Blue)

	envEx := godotenv.Load(".env")
	if envEx != nil {
		utils.Cout("[ERROR] Failed loading environment file: %v", utils.Red, envEx)
	}

	// ==========================================
	ce, cacheEx := strconv.Atoi(os.Getenv("CACHE_EXPIRY"))
	if cacheEx != nil {
		utils.Cout("[ERROR] Could not parse cache expiry time: %v", utils.Red, cacheEx)
	}

	utils.SetupCache(time.Duration(ce) * time.Minute)
	utils.SetupFirebase("firebase-credentials.json")

	// ==========================================
	Mgr, botEx := shards.New("Bot " + os.Getenv("BOT_TOKEN"))
	if botEx != nil {
		utils.Cout("[ERROR] Session creation failed: %v", utils.Red, botEx)
	}

	Mgr.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
		s.UpdateStatusComplex(discordgo.UpdateStatusData{
			Activities: []*discordgo.Activity{{
				Name: "the wind guide my arrows",
				Type: discordgo.ActivityTypeWatching,
				URL:  "",
			}},
		})
	})

	Mgr.AddHandler(func(s *discordgo.Session, e *discordgo.Connect) {
		utils.Cout("[INFO][SHARD-%d] Shard connected.", utils.Cyan, s.ShardID)
	})

	Mgr.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if i.Type == 2 { // ApplicationCommand
			handlers.LinkCommand(s, i)
		} else if i.Type == 3 { // MessageComponent
			handlers.LinkButton(s, i)
		}
	})

	utils.Cout("[INFO] Starting sharding manager.", utils.Cyan)
	shardEx := Mgr.Start()
	if shardEx != nil {
		utils.Cout("[ERROR] Couldn't start the sharding manager: %v", utils.Red, shardEx)
	}

	// ==========================================
	if os.Getenv("REGISTER_COMMANDS") == "1" {
		if os.Getenv("PRODUCTION") == "0" {
			handlers.RegisterCommands(Mgr, os.Getenv("GUILD_ID"))
		} else {
			handlers.RegisterCommands(Mgr, "")
		}
	}

	utils.Cout("\n[SUCCESS] Artemis-Bot is now fully operational.", utils.Green)

	// ==========================================
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	<-stop

	utils.Cout("[INFO] Shutdown sequence initiated.", utils.Blue)
	Mgr.Shutdown()
	utils.Cout("\n[SUCCESS] Shutdown sequence completed.", utils.Green)
}
