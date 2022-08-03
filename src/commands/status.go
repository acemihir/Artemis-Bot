package commands

import (
	"fmt"
	"strings"
	"time"

	"github.com/OnlyF0uR/Artemis-Bot/src/handlers"
	"github.com/OnlyF0uR/Artemis-Bot/src/utils"
	"github.com/bwmarrin/discordgo"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

type dropdownSubmissionData struct {
	SubmissionID      string
	SubmissionMessage *discordgo.Message
	UserID            string
}

var status_dropdown_data = make(map[string]dropdownSubmissionData)

func submissionNotFound(s *discordgo.Session, i *discordgo.InteractionCreate, id string) {
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{
				{
					Title:       handlers.Cfg.Appearance.BotName + " - Status",
					Description: fmt.Sprintf("Submission ``%s`` could not be found", id),
					Color:       utils.WarnEmbedColour,
				},
			},
		},
	})
}

func init() {
	handlers.RegisterCommand(statusCmd)

	handlers.RegisterMessageComponent(statusDropdown)
}

var statusCmd = &handlers.SlashCommand{
	Name:       "status",
	Permission: utils.StaffPermission,
	Exec: func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		id := i.ApplicationCommandData().Options[0].StringValue()

		// Get data from firestore
		res, ex := utils.Firebase.GetFirestore("submissions", id)
		if ex != nil {
			utils.Cout("[ERROR] Get from Firestore failed: %v", utils.Red, ex)
			s.FollowupMessageCreate(s.State.User.ID, i.Interaction, false, &discordgo.WebhookParams{
				Embeds: []*discordgo.MessageEmbed{
					{
						Description: "Oops! A wild error seems to have occured.\n\nPlease try again later, if this error is persistent please report it in our Support discord.",
						Color:       utils.ErrorEmbedColour,
					},
				},
				Components: []discordgo.MessageComponent{
					discordgo.ActionsRow{
						Components: []discordgo.MessageComponent{
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
				Flags: 1 << 6,
			})
			return
		}

		// If there is no submission found with the given ID
		if len(res) == 0 || res["channel_id"] == nil || res["message_id"] == nil {
			submissionNotFound(s, i, id)
			return
		}

		// Check if channel is in guild
		chn, ex := s.Channel(res["channel_id"].(string))
		if ex != nil || chn == nil || chn.GuildID != i.GuildID {
			submissionNotFound(s, i, id)
			return
		}

		// Fetch the message
		if strings.HasPrefix(id, "s_") {
			ex = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Embeds: []*discordgo.MessageEmbed{
						{
							Title:       "Suggestion status",
							Description: fmt.Sprintf("You can change the status of suggestion ``%s`` by selecting a new status in the dropdown menu below.", id),
							Color:       0x336db0,
						},
					},
					Components: []discordgo.MessageComponent{
						discordgo.ActionsRow{
							Components: []discordgo.MessageComponent{
								discordgo.SelectMenu{
									CustomID:    "status_change",
									Placeholder: "Select a status",
									Options: []discordgo.SelectMenuOption{
										{
											Label: "Open",
											Value: "open",
											Emoji: discordgo.ComponentEmoji{
												Name: "⚪",
											},
											Default: true,
										},
										{
											Label: "Considering",
											Value: "considering",
											Emoji: discordgo.ComponentEmoji{
												Name: "🟡",
											},
											Default: false,
										},
										{
											Label: "Approved",
											Value: "approved",
											Emoji: discordgo.ComponentEmoji{
												Name: "🟢",
											},
											Default: false,
										},
										{
											Label: "Rejected",
											Value: "rejected",
											Emoji: discordgo.ComponentEmoji{
												Name: "🔴",
											},
											Default: false,
										},
									},
								},
							},
						},
					},
				},
			})
		} else if strings.HasPrefix(id, "r_") {
			ex = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Embeds: []*discordgo.MessageEmbed{
						{
							Title:       "Report status",
							Description: fmt.Sprintf("You can change the status of report ``%s`` by selecting a new status in the dropdown menu below.", id),
							Color:       0x336db0,
						},
					},
					Components: []discordgo.MessageComponent{
						discordgo.ActionsRow{
							Components: []discordgo.MessageComponent{
								discordgo.SelectMenu{
									CustomID:    "status_change",
									Placeholder: "Select a status",
									Options: []discordgo.SelectMenuOption{
										{
											Label: "Open",
											Value: "open",
											Emoji: discordgo.ComponentEmoji{
												Name: "⚪",
											},
											Default: true,
										},
										{
											Label: "Progressing",
											Value: "progressing",
											Emoji: discordgo.ComponentEmoji{
												Name: "🟡",
											},
											Default: false,
										},
										{
											Label: "Resolved",
											Value: "resolve",
											Emoji: discordgo.ComponentEmoji{
												Name: "🟢",
											},
											Default: false,
										},
									},
								},
							},
						},
					},
				},
			})
		}

		if ex != nil {
			utils.Cout("[ERROR] Could not send message + dropdown: %v", utils.Red, ex)
			utils.ErrorResponse(s, i.Interaction)
			return
		}

		msg, ex := s.InteractionResponse(s.State.User.ID, i.Interaction)
		if ex != nil {
			utils.Cout("[ERROR] Could not get interaction response: %v", utils.Red, ex)
			utils.ErrorResponse(s, i.Interaction)
			return
		}

		sm, ex := s.ChannelMessage(res["channel_id"].(string), res["message_id"].(string))
		if ex != nil {
			submissionNotFound(s, i, id)
			return
		}

		status_dropdown_data[msg.ID] = dropdownSubmissionData{
			SubmissionID:      id,
			UserID:            i.Member.User.ID,
			SubmissionMessage: sm,
		}
	},
}

var statusDropdown = &handlers.MessageComponent{
	ID: "status_change",
	Exec: func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseDeferredMessageUpdate,
		})

		if status_dropdown_data[i.Message.ID] == (dropdownSubmissionData{}) {
			s.FollowupMessageEdit(s.State.User.ID, i.Interaction, i.Message.ID, &discordgo.WebhookEdit{
				Embeds: []*discordgo.MessageEmbed{
					{
						Title:       handlers.Cfg.Appearance.BotName + " - Status",
						Description: "The dropdown menu was expired, please use the status command again.",
						Color:       utils.WarnEmbedColour,
					},
				},
				Components: []discordgo.MessageComponent{},
			})
			return
		}

		if status_dropdown_data[i.Message.ID].UserID != i.Member.User.ID {
			s.FollowupMessageCreate(s.State.User.ID, i.Interaction, false, &discordgo.WebhookParams{
				Embeds: []*discordgo.MessageEmbed{
					{
						Description: "Only the executor of the command can use the dropdown.",
						Color:       utils.ErrorEmbedColour,
					},
				},
				Flags: 1 << 6,
			})
			return
		}

		dropdown_data := status_dropdown_data[i.Message.ID]
		delete(status_dropdown_data, i.Message.ID) // Cleanup

		new_status := i.MessageComponentData().Values[0]

		msg := dropdown_data.SubmissionMessage
		embed := msg.Embeds[0]

		desc_array := strings.Split(embed.Description, "\n")
		desc_array[len(desc_array)-4] = fmt.Sprintf("**Status:** %s (%s)", cases.Title(language.BritishEnglish).String(new_status), time.Now().Format("2006-01-02 15:04:05"))

		colour := utils.PlainEmbedColour
		switch new_status {
		case "approved", "resolved":
			colour = utils.GreenEmbedColour
		case "considering", "progress":
			colour = utils.YellowEmbedColour
		case "reject":
			colour = utils.RedEmbedColour
		}

		embed.Color = colour
		embed.Description = strings.Join(desc_array, "\n")

		_, ex := s.ChannelMessageEditComplex(&discordgo.MessageEdit{
			Channel: dropdown_data.SubmissionMessage.ChannelID,
			ID:      dropdown_data.SubmissionMessage.ID,
			Embeds: []*discordgo.MessageEmbed{
				embed,
			},
		})
		if ex != nil {
			s.FollowupMessageCreate(s.State.User.ID, i.Interaction, false, &discordgo.WebhookParams{
				Embeds: []*discordgo.MessageEmbed{
					{
						Description: "Oops! A wild error seems to have occured.\n\nPlease try again later, if this error is persistent please report it in our Support discord.",
						Color:       utils.ErrorEmbedColour,
					},
				},
				Components: []discordgo.MessageComponent{
					discordgo.ActionsRow{
						Components: []discordgo.MessageComponent{
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
				Flags: 1 << 6,
			})
			return
		}

		s.FollowupMessageEdit(s.State.User.ID, i.Interaction, i.Message.ID, &discordgo.WebhookEdit{
			Embeds: []*discordgo.MessageEmbed{
				{
					Title:       handlers.Cfg.Appearance.BotName + " - Status",
					Description: fmt.Sprintf("You successfully changed the status of ``%s`` to ``%s``.", dropdown_data.SubmissionID, new_status),
					Color:       utils.DefaultEmbedColour,
				},
			},
			Components: []discordgo.MessageComponent{},
		})
	},
}
