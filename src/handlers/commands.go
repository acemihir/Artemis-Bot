package handlers

import (
	"strconv"

	"github.com/bwmarrin/discordgo"
	"github.com/jerskisnow/Artemis-Bot/shards"
	"github.com/jerskisnow/Artemis-Bot/src/commands"
	"github.com/jerskisnow/Artemis-Bot/src/utils"
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
			Name:        "poll",
			Description: "Create & Interact with the polls.",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Name:        "create",
					Description: "Create a personal encrypted note.",
					Type:        discordgo.ApplicationCommandOptionSubCommand,
				},
				{
					Name:        "end",
					Description: "Force end a poll.",
					Type:        discordgo.ApplicationCommandOptionSubCommand,
				},
				{
					Name:        "list",
					Description: "List all active polls.",
					Type:        discordgo.ApplicationCommandOptionSubCommand,
				},
			},
		},
		{
			Name:        "report",
			Description: "Create a report.",
		},
		{
			Name:        "status",
			Description: "Alter the status of a submission.",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Name:        "id",
					Description: "The ID of the submission.",
					Required:    true,
					Type:        discordgo.ApplicationCommandOptionString,
				},
			},
		},
		{
			Name:        "suggest",
			Description: "Create a suggestion.",
		},
	}

	for _, v := range cmds {
		Mgr.ApplicationCommandCreate(guildID, v)
	}
}

func DeleteCommands(Mgr *shards.Manager, guildID string) {
	var s *discordgo.Session

	if guildID == "" {
		// Get the first shard, just getting used for getting the global commands
		s = Mgr.Shards[0].Session
	} else {
		// Get shard for guild
		n, ex := strconv.ParseInt(guildID, 10, 64)
		if ex != nil {
			utils.Cout("[ERROR] Could not parse GuildID to int64: %v", utils.Red, ex)
		}

		s = Mgr.SessionForGuild(n)
	}

	cmds, ex := s.ApplicationCommands(s.State.User.ID, guildID)
	if ex != nil {
		utils.Cout("[ERROR] Could not fetch guild commands: %v", utils.Red, ex)
	}

	for _, v := range cmds {
		Mgr.ApplicationCommandDelete(guildID, v)
	}
}

// We don't need no fancy I/O loops
func LinkCommand(s *discordgo.Session, i *discordgo.InteractionCreate) {
	data := i.ApplicationCommandData()

	switch data.Name {
	case "about":
		commands.AboutCommand(s, i)
	case "config":
		commands.ConfigCommand(s, i)
	case "help":
		commands.HelpCommand(s, i)
	case "notes":
		sbcmd := i.ApplicationCommandData().Options[0].Name
		switch sbcmd {
		case "create":
			commands.NotesCreateCommand(s, i)
		case "delete":
			commands.NotesDeleteCommand(s, i)
		case "list":
			commands.NotesListCommand(s, i)
		}
	case "poll":
		sbcmd := i.ApplicationCommandData().Options[0].Name
		switch sbcmd {
		case "create":
			commands.PollCreateCommand(s, i)
		case "end":
			commands.PollEndCommand(s, i)
		case "list":
			commands.PollListCommand(s, i)
		}
	case "report":
		commands.ReportCommand(s, i)
	case "status":
		commands.StatusCommand(s, i)
	case "suggest":
		commands.SuggestCommand(s, i)
	}
}
