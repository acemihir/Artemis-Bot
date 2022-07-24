package commands

import (
	"github.com/OnlyF0uR/Artemis-Bot/src/handlers"
	"github.com/bwmarrin/discordgo"
)

func init() {
	handlers.RegisterCommand(aboutCmd)
}

var aboutCmd = &handlers.SlashCommand{
	Name: "about",
	Exec: func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Embeds: []*discordgo.MessageEmbed{
					{
						Title:       "About Artemis",
						Description: "[Artemis-Bot](https://github.com/OnlyF0uR/Artemis-Bot), derived from the Greek goddess of the hunt, is an ambitious discord bot project. It mainly helps establish a perfect connection between management and the community. The bot could be characterized as a multi-function/purpose bot that fits all your needs when it comes to managing and interacting with your community.",
						Color:       0x336db0,
					},
				},
				Components: []discordgo.MessageComponent{
					discordgo.ActionsRow{
						Components: []discordgo.MessageComponent{
							discordgo.Button{
								Emoji: discordgo.ComponentEmoji{
									Name: "🤖",
								},
								Label: "Invite",
								Style: discordgo.LinkButton,
								URL:   "https://top.gg/bot/566616056165302282/invite/",
							},
							discordgo.Button{
								Emoji: discordgo.ComponentEmoji{
									Name: "📰",
								},
								Label: "Vote",
								Style: discordgo.LinkButton,
								URL:   "https://top.gg/bot/566616056165302282/vote",
							},
							discordgo.Button{
								Emoji: discordgo.ComponentEmoji{
									Name: "💰",
								},
								Label: "Contribute",
								Style: discordgo.LinkButton,
								URL:   "https://github.com/OnlyF0uR/Artemis-Bot/wiki/Donating",
							},
							discordgo.Button{
								Emoji: discordgo.ComponentEmoji{
									Name: "👥",
								},
								Label: "Support",
								Style: discordgo.LinkButton,
								URL:   "https://discord.gg/3SYg3M5",
							},
						},
					},
				},
			},
		})
	},
}
