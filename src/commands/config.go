package commands

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/jerskisnow/Artemis-Bot/src/utils"
)

// ===========================================
// MENU FUNCTIONS
// ===========================================
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

// ===========================================
// INDIVIDUAL FUNCTION BUTTON RESPONSES
// ===========================================
func ConfigAuthStaffroleButton(s *discordgo.Session, i *discordgo.InteractionCreate) {
	ex := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseModal,
		Data: &discordgo.InteractionResponseData{
			CustomID: "modals_config_auth_staffrole",
			Title:    "Enter the staffrole",
			Components: []discordgo.MessageComponent{
				discordgo.ActionsRow{
					Components: []discordgo.MessageComponent{
						discordgo.TextInput{
							CustomID:    "staffrole",
							Label:       "Name or ID of role",
							Style:       discordgo.TextInputShort,
							Placeholder: "@moderator",
							Required:    true,
							MaxLength:   300,
							MinLength:   3,
						},
					},
				},
			},
		},
	})
	if ex != nil {
		utils.Cout("[ERROR] Could not open up the modal: %v", utils.Red, ex)
		utils.ErrorResponse(s, i.Interaction)
		return
	}
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

// ===========================================
// REACTIONS TO THE SUBMISSIONS OF THE MODALS
// ===========================================
func ConfigAuthStaffroleModal(s *discordgo.Session, i *discordgo.InteractionCreate) {
	role_id := i.ModalSubmitData().Components[0].(*discordgo.ActionsRow).Components[0].(*discordgo.TextInput).Value

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
	})

	fmt.Println(role_id)
	var role *discordgo.Role = nil

	// Get all the roles
	roles, ex := s.GuildRoles(i.GuildID)
	if ex != nil {
		utils.Cout("[ERROR] Could not get server roles: %v", utils.Red, ex)
		utils.ErrorResponse(s, i.Interaction)
		return
	}

	// Check if number
	if _, err := strconv.Atoi(role_id); err == nil {
		// Get the role by id
		for _, r := range roles {
			if r.ID == role_id {
				role = r
				break
			}
		}
	} else {
		// Remove the @ that is possibly in front
		tmp := strings.Replace(role_id, "@", "", 1)

		// Get the role by name
		for _, r := range roles {
			if r.Name == tmp {
				role = r
				break
			}
		}
	}

	// If the role was not found then notify the user
	if role == nil {
		s.FollowupMessageCreate(s.State.User.ID, i.Interaction, false, &discordgo.WebhookParams{
			Embeds: []*discordgo.MessageEmbed{
				{
					Title:       "Artemis - Config",
					Description: "The role you tried to configure was not found.",
					Color:       utils.WarnEmbedColour,
				},
			},
		})
		return
	}

	// Save the data to the database
	ex = utils.Firebase.SetFirestore("guilds", i.GuildID, map[string]interface{}{
		"staffrole_id": role.ID,
	}, true)
	if ex != nil {
		utils.Cout("[ERROR] Could not set the staffrole: %v", utils.Red, ex)
		utils.ErrorResponse(s, i.Interaction)
		return
	}

	// Send a message stating that the role was set
	s.FollowupMessageCreate(s.State.User.ID, i.Interaction, false, &discordgo.WebhookParams{
		Embeds: []*discordgo.MessageEmbed{
			{
				Title:       "Artemis - Config",
				Description: "The stafffrole is now set to ``" + role.Name + " (" + role.ID + ")``.",
				Color:       utils.DefaultEmbedColour,
			},
		},
	})
}
