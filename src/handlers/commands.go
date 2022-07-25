package handlers

import (
	"strconv"

	"github.com/OnlyF0uR/Artemis-Bot/shards"
	"github.com/OnlyF0uR/Artemis-Bot/src/utils"
	"github.com/bwmarrin/discordgo"
)

// GuildID should be empty in production
func SubmitCommands(Mgr *shards.Manager, guildID string) {
	cmds := []*discordgo.ApplicationCommand{
		{
			Name:        "about",
			Description: "Obtain information about the bot.",
		},
		{
			Name:        "config",
			Description: "Configure the bot to fit your needs. (ADMINISTRATOR)",
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
							Name:  "Note",
							Value: "note",
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
			Name:        "note",
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
					Options: []*discordgo.ApplicationCommandOption{
						{
							Name:        "title",
							Description: "The title of the note.",
							Required:    true,
							Type:        discordgo.ApplicationCommandOptionString,
						},
					},
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
					Options: []*discordgo.ApplicationCommandOption{
						{
							Name:        "amount",
							Description: "The amount of options for the poll.",
							Required:    true,
							Type:        discordgo.ApplicationCommandOptionString,
						},
					},
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

func RetractCommands(Mgr *shards.Manager, guildID string) {
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

// ======================
// COMMAND HANDLER
// ======================
var cmds = map[string]*SlashCommand{}

type SlashCommand struct {
	Name       string
	Permission int64
	Exec       func(*discordgo.Session, *discordgo.InteractionCreate)
}

func RegisterCommand(cmd *SlashCommand) {
	cmds[cmd.Name] = cmd
}

func LinkCommand(s *discordgo.Session, i *discordgo.InteractionCreate) {
	data := i.ApplicationCommandData()

	// Check if command exists
	if v, ok := cmds[data.Name]; ok {
		if v.Permission != 0 {
			// Check if the user has that permission
			if !utils.HasPermission(i.Member.Permissions, v.Permission) {
				// Insufficient permission
				s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Embeds: []*discordgo.MessageEmbed{
							{
								Description: "You do not have permission to use this command.",
								Color:       utils.WarnEmbedColour,
							},
						},
					},
				})
				return
			}
		}

		v.Exec(s, i)
	}
}
