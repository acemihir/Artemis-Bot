package commands

import (
	"github.com/bwmarrin/discordgo"
	"github.com/jerskisnow/Artemis-Bot/src/handlers"
	"github.com/jerskisnow/Artemis-Bot/src/utils"
)

func init() {
	handlers.RegisterCommand(pollCommand)
}

var pollCommand = &handlers.SlashCommand{
	Name: "poll",
	Exec: func(s *discordgo.Session, ic *discordgo.InteractionCreate) {
		sbcmd := i.ApplicationCommandData().Options[0].Name
		switch sbcmd {
		case "create":
			if !utils.HasPermission(i.Member.Permissions, utils.StaffPermission) {
				s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Embeds: []*discordgo.MessageEmbed{
							{
								Description: "You do not have permission to use this command. (``MANAGE_MESSAGES``)",
								Color:       utils.WarnEmbedColour,
							},
						},
					},
				})
				return
			}
			commands.pollCreateSubcmd(s, i)
		case "end":
			if !utils.HasPermission(i.Member.Permissions, utils.StaffPermission) {
				s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Embeds: []*discordgo.MessageEmbed{
							{
								Description: "You do not have permission to use this command. (``MANAGE_MESSAGES``)",
								Color:       utils.WarnEmbedColour,
							},
						},
					},
				})
				return
			}
			commands.pollEndSubcmd(s, i)
		case "list":
			commands.pollListSubcmd(s, i)
		}
	},
}

func pollCreateSubcmd(s *discordgo.Session, i *discordgo.InteractionCreate) {
	// TODO: This
}

func pollEndSubcmd(s *discordgo.Session, i *discordgo.InteractionCreate) {
	// TODO: This
}

func pollListSubcmd(s *discordgo.Session, i *discordgo.InteractionCreate) {
	// TODO: This
}
