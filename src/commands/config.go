package commands

import (
	"strconv"
	"strings"

	"github.com/OnlyF0uR/Artemis-Bot/src/handlers"
	"github.com/OnlyF0uR/Artemis-Bot/src/utils"
	"github.com/bwmarrin/discordgo"
)

func init() {
	handlers.RegisterCommand(configCommand)

	handlers.RegisterMessageComponent(configMainAuthButton)
	handlers.RegisterMessageComponent(configMainChannelsButton)
	handlers.RegisterMessageComponent(configAuthStaffroleButton)
	handlers.RegisterMessageComponent(configChnsSugButton)
	handlers.RegisterMessageComponent(configChnsRepButton)

	handlers.RegisterModal(configAuthStaffroleModal)
	handlers.RegisterModal(configChnsSugModal)
	handlers.RegisterModal(configChnsRepModal)
}

// ===========================================
// MENU FUNCTIONS
// ===========================================
var configCommand = &handlers.SlashCommand{
	Name:       "config",
	Permission: utils.AdminPermission,
	Exec: func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Embeds: []*discordgo.MessageEmbed{
					{
						Title:       "Config - Menu",
						Description: "This is the config menu for " + handlers.Cfg.Appearance.BotName + ", from here you can simply select a category and follow further instructions.",
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
						},
					},
				},
			},
		})
	},
}

// ===========================================
// MAIN MENUS
// ===========================================
var configMainAuthButton = &handlers.MessageComponent{
	ID: "cfg_main_auth",
	Exec: func(s *discordgo.Session, i *discordgo.InteractionCreate) {
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
	},
}

var configMainChannelsButton = &handlers.MessageComponent{
	ID: "cfg_main_chns",
	Exec: func(s *discordgo.Session, i *discordgo.InteractionCreate) {
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
	},
}

// ===========================================
// SECONDARY MENUS
// ===========================================
var configAuthStaffroleButton = &handlers.MessageComponent{
	ID: "cfg_auth_staffrole",
	Exec: func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		ex := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseModal,
			Data: &discordgo.InteractionResponseData{
				CustomID: "config_auth_staffrole",
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
								MaxLength:   30,
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
	},
}

var configChnsSugButton = &handlers.MessageComponent{
	ID: "cfg_chns_sug",
	Exec: func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		ex := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseModal,
			Data: &discordgo.InteractionResponseData{
				CustomID: "config_channels_sugchannel",
				Title:    "Enter a channel",
				Components: []discordgo.MessageComponent{
					discordgo.ActionsRow{
						Components: []discordgo.MessageComponent{
							discordgo.TextInput{
								CustomID:    "sugchannel",
								Label:       "Name or ID of channel",
								Style:       discordgo.TextInputShort,
								Placeholder: "#suggestions",
								Required:    true,
								MaxLength:   30,
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
	},
}

var configChnsRepButton = &handlers.MessageComponent{
	ID: "cfg_chns_rep",
	Exec: func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		ex := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseModal,
			Data: &discordgo.InteractionResponseData{
				CustomID: "config_channels_repchannel",
				Title:    "Enter a channel",
				Components: []discordgo.MessageComponent{
					discordgo.ActionsRow{
						Components: []discordgo.MessageComponent{
							discordgo.TextInput{
								CustomID:    "repchannel",
								Label:       "Name or ID of channel",
								Style:       discordgo.TextInputShort,
								Placeholder: "#reports",
								Required:    true,
								MaxLength:   30,
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
	},
}

// ===========================================
// SUBMISSION REACTION
// ===========================================
var configAuthStaffroleModal = &handlers.Modal{
	ID: "config_auth_staffrole",
	Exec: func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		role_id := i.ModalSubmitData().Components[0].(*discordgo.ActionsRow).Components[0].(*discordgo.TextInput).Value

		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
		})

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
						Title:       handlers.Cfg.Appearance.BotName + " - Config",
						Description: "The role you tried to configure was not found.",
						Color:       utils.WarnEmbedColour,
					},
				},
			})
			return
		}

		// Save the data to the database
		ex = utils.Firebase.SetFirestore("guilds", i.GuildID, map[string]interface{}{
			"staffrole": role.ID,
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
					Title:       handlers.Cfg.Appearance.BotName + " - Config",
					Description: "The stafffrole is now set to ``" + role.Name + " (" + role.ID + ")``.",
					Color:       utils.DefaultEmbedColour,
				},
			},
		})
	},
}

var configChnsSugModal = &handlers.Modal{
	ID: "config_channels_sugchannel",
	Exec: func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		chn_id := i.ModalSubmitData().Components[0].(*discordgo.ActionsRow).Components[0].(*discordgo.TextInput).Value

		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
		})

		var chn *discordgo.Channel = nil

		// Get all the channels
		channels, ex := s.GuildChannels(i.GuildID)
		if ex != nil {
			utils.Cout("[ERROR] Could not get server channels: %v", utils.Red, ex)
			utils.ErrorResponse(s, i.Interaction)
			return
		}

		// Check if number
		if _, err := strconv.Atoi(chn_id); err == nil {
			// Get the channel by id
			for _, c := range channels {
				if c.ID == chn_id {
					chn = c
					break
				}
			}
		} else {
			tmp := strings.Replace(chn_id, "#", "", 1)

			// Get the channel by name
			for _, r := range channels {
				if r.Name == tmp {
					chn = r
					break
				}
			}
		}

		// If the channel was not found then notify the user
		if chn == nil {
			s.FollowupMessageCreate(s.State.User.ID, i.Interaction, false, &discordgo.WebhookParams{
				Embeds: []*discordgo.MessageEmbed{
					{
						Title:       handlers.Cfg.Appearance.BotName + " - Config",
						Description: "The channel you tried to configure was not found.",
						Color:       utils.WarnEmbedColour,
					},
				},
			})
			return
		}

		// Save the data to the database
		ex = utils.Firebase.SetFirestore("guilds", i.GuildID, map[string]interface{}{
			"sug_channel": chn.ID,
		}, true)
		if ex != nil {
			utils.Cout("[ERROR] Could not set the sugchannel: %v", utils.Red, ex)
			utils.ErrorResponse(s, i.Interaction)
			return
		}

		// Send a message stating that the channel was set
		s.FollowupMessageCreate(s.State.User.ID, i.Interaction, false, &discordgo.WebhookParams{
			Embeds: []*discordgo.MessageEmbed{
				{
					Title:       handlers.Cfg.Appearance.BotName + " - Config",
					Description: "The suggestions channel is now set to ``" + chn.Name + " (" + chn.ID + ")``.",
					Color:       utils.DefaultEmbedColour,
				},
			},
		})
	},
}

var configChnsRepModal = &handlers.Modal{
	ID: "config_channels_repchannel",
	Exec: func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		chn_id := i.ModalSubmitData().Components[0].(*discordgo.ActionsRow).Components[0].(*discordgo.TextInput).Value

		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
		})

		var chn *discordgo.Channel = nil

		// Get all the channels
		channels, ex := s.GuildChannels(i.GuildID)
		if ex != nil {
			utils.Cout("[ERROR] Could not get server channels: %v", utils.Red, ex)
			utils.ErrorResponse(s, i.Interaction)
			return
		}

		// Check if number
		if _, err := strconv.Atoi(chn_id); err == nil {
			// Get the channel by id
			for _, c := range channels {
				if c.ID == chn_id {
					chn = c
					break
				}
			}
		} else {
			tmp := strings.Replace(chn_id, "#", "", 1)

			// Get the channel by name
			for _, r := range channels {
				if r.Name == tmp {
					chn = r
					break
				}
			}
		}

		// If the channel was not found then notify the user
		if chn == nil {
			s.FollowupMessageCreate(s.State.User.ID, i.Interaction, false, &discordgo.WebhookParams{
				Embeds: []*discordgo.MessageEmbed{
					{
						Title:       handlers.Cfg.Appearance.BotName + " - Config",
						Description: "The channel you tried to configure was not found.",
						Color:       utils.WarnEmbedColour,
					},
				},
			})
			return
		}

		// Save the data to the database
		ex = utils.Firebase.SetFirestore("guilds", i.GuildID, map[string]interface{}{
			"rep_channel": chn.ID,
		}, true)
		if ex != nil {
			utils.Cout("[ERROR] Could not set the sugchannel: %v", utils.Red, ex)
			utils.ErrorResponse(s, i.Interaction)
			return
		}

		// Send a message stating that the channel was set
		s.FollowupMessageCreate(s.State.User.ID, i.Interaction, false, &discordgo.WebhookParams{
			Embeds: []*discordgo.MessageEmbed{
				{
					Title:       handlers.Cfg.Appearance.BotName + " - Config",
					Description: "The reports channel is now set to ``" + chn.Name + " (" + chn.ID + ")``.",
					Color:       utils.DefaultEmbedColour,
				},
			},
		})
	},
}
