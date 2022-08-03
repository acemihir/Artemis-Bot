package main

import (
	"os"
	"os/signal"
	"time"

	"github.com/OnlyF0uR/Artemis-Bot/shards"
	"github.com/bwmarrin/discordgo"

	// Activate all the init functions
	_ "github.com/OnlyF0uR/Artemis-Bot/src/commands"

	"github.com/OnlyF0uR/Artemis-Bot/src/handlers"
	"github.com/OnlyF0uR/Artemis-Bot/src/utils"
)

var (
	Mgr *shards.Manager
)

func main() {
	// Load the config
	handlers.LoadConfig()

	utils.Cout("[INFO] Start sequence initiated.\n", utils.Blue)
	// ==========================================

	utils.SetupCache(time.Duration(handlers.Cfg.Data.CacheExpiry) * time.Minute)
	utils.SetupFirebase("firebase-credentials.json")

	handlers.RegisterTasks(handlers.Cfg.AppMode == "production")

	// ==========================================
	Mgr, botEx := shards.New("Bot " + handlers.Cfg.Client.AuthToken)
	if botEx != nil {
		utils.Cout("[ERROR] Session creation failed: %v", utils.Red, botEx)
		os.Exit(1)
	}

	Mgr.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
		s.UpdateStatusComplex(discordgo.UpdateStatusData{
			Activities: []*discordgo.Activity{{
				Name: handlers.Cfg.Client.ActivityText,
				Type: discordgo.ActivityType(handlers.Cfg.Client.ActivityType),
				URL:  handlers.Cfg.Client.ActivityUrl,
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
	if handlers.Cfg.AppMode == "production" {
		if handlers.Cfg.Commands.RetractAll {
			// Delete globally
			handlers.RetractCommands(Mgr, "")
		}

		if handlers.Cfg.Commands.SubmitAll {
			// Register globally
			handlers.SubmitCommands(Mgr, "")
		}
	} else {
		if handlers.Cfg.Commands.RetractAll {
			// Delete for development guild
			handlers.RetractCommands(Mgr, handlers.Cfg.Client.GuildID)
		}

		if handlers.Cfg.Commands.SubmitAll {
			// Register for development guild
			handlers.SubmitCommands(Mgr, handlers.Cfg.Client.GuildID)
		}
	}

	utils.Cout("\n[SUCCESS] "+handlers.Cfg.Appearance.BotName+"-Bot is now fully operational.", utils.Green)

	// ==========================================
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	<-stop

	utils.Cout("[INFO] Shutdown sequence initiated.", utils.Blue)
	Mgr.Shutdown()
	if handlers.Cfg.AppMode == "production" {
		handlers.ShutdownTasks()
	}
	utils.Cout("\n[SUCCESS] Shutdown sequence completed.", utils.Green)
}
