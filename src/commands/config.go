package commands

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

func ConfigCommand(s *discordgo.Session, i *discordgo.InteractionCreate) {
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{
				{
					Title:       "Config - Menu",
					Description: "This is the config menu for Artemis, from here you can simply select a category and follow further instructions.",
					Color:       0x336db0,
				},
			},
			Components: []discordgo.MessageComponent{
				discordgo.ActionsRow{
					Components: []discordgo.MessageComponent{
						discordgo.Button{
							CustomID: "cfg_main_auth",
							Emoji: discordgo.ComponentEmoji{
								Name: "👥",
							},
							Label: "Auth",
							Style: discordgo.PrimaryButton,
						},
						discordgo.Button{
							CustomID: "cfg_main_chns",
							Emoji: discordgo.ComponentEmoji{
								Name: "#️⃣",
							},
							Label: "Channels",
							Style: discordgo.PrimaryButton,
						},
						discordgo.Button{
							CustomID: "cfg_main_appear",
							Emoji: discordgo.ComponentEmoji{
								Name: "👗",
							},
							Label: "Appearance",
							Style: discordgo.PrimaryButton,
						},
						// discordgo.Button{
						// 	CustomID: "cfg_main_misc",
						// 	Emoji: discordgo.ComponentEmoji{
						// 		Name: "🛰️",
						// 	},
						// 	Label: "Misc",
						// 	Style: discordgo.PrimaryButton,
						// },
					},
				},
			},
		},
	})
}

func ConfigMainAuthButton(s *discordgo.Session, i *discordgo.InteractionCreate) {
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseUpdateMessage,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{
				{
					Title:       "Config - Auth",
					Description: "Settings related to authorization can be configured here, this includes staffroles etc.",
					Color:       0x336db0,
				},
			},
			Components: []discordgo.MessageComponent{
				discordgo.ActionsRow{
					Components: []discordgo.MessageComponent{
						discordgo.Button{
							CustomID: "cfg_auth_staffrole",
							Emoji: discordgo.ComponentEmoji{
								Name: "🦸",
							},
							Label: "Staffrole",
							Style: discordgo.SecondaryButton,
						},
					},
				},
			},
		},
	})
}

func ConfigMainChannelsButton(s *discordgo.Session, i *discordgo.InteractionCreate) {
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseUpdateMessage,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{
				{
					Title:       "Config - Channels",
					Description: "All the different channels can be set in this section.",
					Color:       0x336db0,
				},
			},
			Components: []discordgo.MessageComponent{
				discordgo.ActionsRow{
					Components: []discordgo.MessageComponent{
						discordgo.Button{
							CustomID: "cfg_chns_sug",
							Emoji: discordgo.ComponentEmoji{
								Name: "❔",
							},
							Label: "Suggestions",
							Style: discordgo.SecondaryButton,
						},
						discordgo.Button{
							CustomID: "cfg_chns_rep",
							Emoji: discordgo.ComponentEmoji{
								Name: "❗",
							},
							Label: "Reports",
							Style: discordgo.SecondaryButton,
						},
					},
				},
			},
		},
	})
}

func ConfigMainAppearanceButton(s *discordgo.Session, i *discordgo.InteractionCreate) {
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseUpdateMessage,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{
				{
					Title:       "Config - Appearance",
					Description: "This section of the configuration is designated for altering the appearance of the bot such as the emojis the bot uses.",
					Color:       0x336db0,
				},
			},
			Components: []discordgo.MessageComponent{
				discordgo.ActionsRow{
					Components: []discordgo.MessageComponent{
						discordgo.Button{
							CustomID: "cfg_appear_upvote",
							Emoji: discordgo.ComponentEmoji{
								Name: "⬆️",
							},
							Label: "Upvote Emoji",
							Style: discordgo.SecondaryButton,
						},
						discordgo.Button{
							CustomID: "cfg_appear_downvote",
							Emoji: discordgo.ComponentEmoji{
								Name: "⬇️",
							},
							Label: "Downvote Emoji",
							Style: discordgo.SecondaryButton,
						},
					},
				},
			},
		},
	})
}

// func ConfigMainMiscButton(s *discordgo.Session, i *discordgo.InteractionCreate) {}

func ConfigAuthStaffroleButton(s *discordgo.Session, i *discordgo.InteractionCreate) {
	fmt.Println("TODO")
	// TODO: This
}

func ConfigChnsSugButton(s *discordgo.Session, i *discordgo.InteractionCreate) {
	fmt.Println("TODO")
	// TODO: This
}

func ConfigChnsRepButton(s *discordgo.Session, i *discordgo.InteractionCreate) {
	fmt.Println("TODO")
	// TODO: This
}

func ConfigAppearUpvoteButton(s *discordgo.Session, i *discordgo.InteractionCreate) {
	fmt.Println("TODO")
	// TODO: This
}

func ConfigAppearDownvoteButton(s *discordgo.Session, i *discordgo.InteractionCreate) {
	fmt.Println("TODO")
	// TODO: This
}
