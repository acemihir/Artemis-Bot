package handlers

import (
	"log"

	"github.com/bwmarrin/discordgo"
	"github.com/jerskisnow/Suggestions/src/commands"
)

// GuildID should be empty in production
func RegisterCommands(s *discordgo.Session, guildID string) {
	cmds := []*discordgo.ApplicationCommand{
		{
			Name:        "about",
			Description: "Obtain information about the bot",
		},
	}

	log.Println("Creating commands.")

	for _, v := range cmds {
		_, ex := s.ApplicationCommandCreate(s.State.User.ID, guildID, v)
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
	} else if data.Name == "suggest" {
		commands.ExecSuggest(s, i)
	}
}
