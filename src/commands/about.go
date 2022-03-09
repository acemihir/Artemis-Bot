package commands

import "github.com/bwmarrin/discordgo"

func ExecAbout(s *discordgo.Session, i *discordgo.InteractionCreate) {
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{
				&discordgo.MessageEmbed{
					Title:       "About - Artemis",
					Description: "Artemis is an ambitious discord bot project. It mainly helps establish a perfect connection between management and the community.",
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
							URL:   "https://placeholder.com/",
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
}
