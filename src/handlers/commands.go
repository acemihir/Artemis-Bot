package handlers

import (
	"log"

	"github.com/bwmarrin/discordgo"
	"github.com/jerskisnow/Suggestions/src/commands"
	"github.com/servusdei2018/shards"
)

// GuildID should be empty in production
func RegisterCommands(Mgr *shards.Manager, guildID string) {
	cmds := []*discordgo.ApplicationCommand{
		{
			Name:        "about",
			Description: "Obtain information about the bot.",
		},
		{
			Name:        "config",
			Description: "Configure the bot to fit your needs. (MANAGE_MESSAGES)",
		},
		{
			Name:        "help",
			Description: "Receive a list of all commands.",
		},
		{
			Name:        "report",
			Description: "Create a report.",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Name:        "description",
					Description: "A brief description of your report.",
					Required:    true,
					Type:        discordgo.ApplicationCommandOptionString,
				},
			},
		},
		{
			Name:        "suggest",
			Description: "Create a suggestion.",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Name:        "description",
					Description: "A brief description of your suggestion.",
					Required:    true,
					Type:        discordgo.ApplicationCommandOptionString,
				},
			},
		},
	}

	log.Println("Creating commands.")

	for _, v := range cmds {
		ex := Mgr.ApplicationCommandCreate(guildID, v)
		if ex != nil {
			log.Panicf("Cannot create '%v' command: %v", v.Name, ex)
		}
	}
}

// We don't need no fancy I/O loops
func LinkCommand(s *discordgo.Session, i *discordgo.InteractionCreate) {
	data := i.ApplicationCommandData()

	if data.Name == "about" {
		commands.ExecAbout(s, i)
	} else if data.Name == "config" {
		commands.ExecConfig(s, i)
	} else if data.Name == "help" {
		commands.ExecHelp(s, i)
	} else if data.Name == "notes" {
		commands.ExecNotes(s, i)
	} else if data.Name == "poll" {
		commands.ExecPoll(s, i)
	} else if data.Name == "report" {
		commands.ExecReport(s, i)
	} else if data.Name == "setup" {
		commands.ExecReport(s, i)
	} else if data.Name == "status" {
		commands.ExecStatus(s, i)
	} else if data.Name == "suggest" {
		commands.ExecSuggest(s, i)
	}
}
