package commands

import (
	"github.com/OnlyF0uR/Artemis-Bot/src/handlers"
	"github.com/OnlyF0uR/Artemis-Bot/src/utils"
	"github.com/bwmarrin/discordgo"
)

func init() {
	handlers.RegisterCommand(pollCommand)
}

var pollCommand = &handlers.SlashCommand{
	Name: "poll",
	Exec: func(s *discordgo.Session, i *discordgo.InteractionCreate) {
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
			pollCreateSubcmd(s, i)
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
			pollEndSubcmd(s, i)
		case "list":
			pollListSubcmd(s, i)
		}
	},
}

func pollCreateSubcmd(s *discordgo.Session, i *discordgo.InteractionCreate) {
	// TODO: This
	utils.ComingSoonResponse(s, i.Interaction)
}

func pollEndSubcmd(s *discordgo.Session, i *discordgo.InteractionCreate) {
	// TODO: This
	utils.ComingSoonResponse(s, i.Interaction)
}

func pollListSubcmd(s *discordgo.Session, i *discordgo.InteractionCreate) {
	// TODO: This
	utils.ComingSoonResponse(s, i.Interaction)
}
