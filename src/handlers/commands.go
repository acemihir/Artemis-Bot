package handlers

import (
	"log"
	"strconv"

	"github.com/bwmarrin/discordgo"
	"github.com/jerskisnow/Suggestions/shards"
	"github.com/jerskisnow/Suggestions/src/commands"
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
			Options: []*discordgo.ApplicationCommandOption{
				{
					Name:        "command",
					Description: "The command you would like to read more about.",
					Required:    false,
					Type:        discordgo.ApplicationCommandOptionString,
					Choices: []*discordgo.ApplicationCommandOptionChoice{
						{
							Name:  "About",
							Value: "about",
						},
						{
							Name:  "Config",
							Value: "config",
						},
						{
							Name:  "Help",
							Value: "help",
						},
						{
							Name:  "Notes",
							Value: "notes",
						},
						{
							Name:  "Poll",
							Value: "poll",
						},
						{
							Name:  "Report",
							Value: "report",
						},
						{
							Name:  "Status",
							Value: "status",
						},
						{
							Name:  "Suggest",
							Value: "suggest",
						},
					},
				},
			},
		},
		{
			Name:        "notes",
			Description: "Create & Interact with your own notes.",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Name:        "create",
					Description: "Create a personal encrypted note.",
					Type:        discordgo.ApplicationCommandOptionSubCommand,
				},
				{
					Name:        "delete",
					Description: "Delete a personal note.",
					Type:        discordgo.ApplicationCommandOptionSubCommand,
				},
				{
					Name:        "list",
					Description: "Obtain a list with all your personal notes.",
					Type:        discordgo.ApplicationCommandOptionSubCommand,
				},
			},
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

	log.Println("[INFO] Creating commands.")

	for _, v := range cmds {
		Mgr.ApplicationCommandCreate(guildID, v)
	}
	log.Println("[INFO] Finished registering all commands.")
}

func DeleteCommands(Mgr *shards.Manager, guildID string) {
	n, ex := strconv.ParseInt(guildID, 10, 64)
	if ex != nil {
		log.Fatalf("[ERROR] Could not parse GuildID to int64. (DeleteCommands)")
	}

	s := Mgr.SessionForGuild(n)
	cmds, ex := s.ApplicationCommands(s.State.User.ID, guildID)
	if ex != nil {
		log.Fatalf("[ERROR] Could not fetch guild commands. (DeleteCommands)")
	}

	for _, v := range cmds {
		Mgr.ApplicationCommandDelete(guildID, v)
	}
}

// We don't need no fancy I/O loops
func LinkCommand(s *discordgo.Session, i *discordgo.InteractionCreate) {
	data := i.ApplicationCommandData()

	if data.Name == "about" {
		commands.AboutCommand(s, i)
	} else if data.Name == "config" {
		commands.ConfigCommand(s, i)
	} else if data.Name == "help" {
		commands.HelpCommand(s, i)
	} else if data.Name == "notes" {
		commands.NotesCommand(s, i)
	} else if data.Name == "poll" {
		commands.PollCommand(s, i)
	} else if data.Name == "report" {
		commands.ReportCommand(s, i)
	} else if data.Name == "status" {
		commands.StatusCommand(s, i)
	} else if data.Name == "suggest" {
		commands.SuggestCommand(s, i)
	}
}
