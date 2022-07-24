package main

import (
	"os"
	"os/signal"
	"strconv"
	"time"

	"github.com/OnlyF0uR/Artemis-Bot/shards"
	"github.com/bwmarrin/discordgo"

	// Activate all the init functions
	_ "github.com/OnlyF0uR/Artemis-Bot/src/commands"

	"github.com/OnlyF0uR/Artemis-Bot/src/handlers"
	"github.com/OnlyF0uR/Artemis-Bot/src/utils"
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
		os.Exit(1)
	}

	// ==========================================
	ce, cacheEx := strconv.Atoi(os.Getenv("CACHE_EXPIRY"))
	if cacheEx != nil {
		utils.Cout("[ERROR] Could not parse cache expiry time: %v", utils.Red, cacheEx)
		os.Exit(1)
	}

	utils.SetupCache(time.Duration(ce) * time.Minute)
	utils.SetupFirebase("firebase-credentials.json")
	handlers.RegisterTasks()

	// ==========================================
	Mgr, botEx := shards.New("Bot " + os.Getenv("BOT_TOKEN"))
	if botEx != nil {
		utils.Cout("[ERROR] Session creation failed: %v", utils.Red, botEx)
		os.Exit(1)
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
		if i.Type == discordgo.InteractionApplicationCommand {
			handlers.LinkCommand(s, i)
		} else if i.Type == discordgo.InteractionMessageComponent {
			handlers.LinkMessageComponent(s, i)
		} else if i.Type == discordgo.InteractionModalSubmit {
			handlers.LinkModal(s, i)
		}
	})

	utils.Cout("[INFO] Starting sharding manager.", utils.Cyan)
	shardEx := Mgr.Start()
	if shardEx != nil {
		utils.Cout("[ERROR] Couldn't start the sharding manager: %v", utils.Red, shardEx)
		os.Exit(1)
	}

	// ==========================================
	if os.Getenv("DELETE_COMMANDS") == "1" {
		if os.Getenv("PRODUCTION") == "0" {
			handlers.RetractCommands(Mgr, os.Getenv("GUILD_ID"))
		} else {
			handlers.RetractCommands(Mgr, "")
		}
	}
	if os.Getenv("REGISTER_COMMANDS") == "1" {
		if os.Getenv("PRODUCTION") == "0" {
			handlers.SubmitCommands(Mgr, os.Getenv("GUILD_ID"))
		} else {
			handlers.SubmitCommands(Mgr, "")
		}
	}

	utils.Cout("\n[SUCCESS] Artemis-Bot is now fully operational.", utils.Green)

	// ==========================================
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	<-stop

	utils.Cout("[INFO] Shutdown sequence initiated.", utils.Blue)
	Mgr.Shutdown()
	if os.Getenv("PRODUCTION") == "1" {
		handlers.ShutdownTasks()
	}
	utils.Cout("\n[SUCCESS] Shutdown sequence completed.", utils.Green)
}