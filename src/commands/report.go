package commands

import (
	"fmt"
	"strings"

	"github.com/OnlyF0uR/Artemis-Bot/src/handlers"
	"github.com/OnlyF0uR/Artemis-Bot/src/utils"
	"github.com/bwmarrin/discordgo"
)

var rep_modal_data = make(map[string]string)

func init() {
	handlers.RegisterCommand(reportCmd)

	handlers.RegisterModal(reportCreateModal)
}

var reportCmd = &handlers.SlashCommand{
	Name: "report",
	Exec: func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		res, ex := utils.Firebase.GetFirestore("guilds", i.GuildID)
		if ex != nil {
			utils.Cout("[ERROR] Get from Firestore failed: %v", utils.Red, ex)
			utils.ErrorResponse(s, i.Interaction)
			return
		}

		if len(res) == 0 || res["rep_channel"] == nil {
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Embeds: []*discordgo.MessageEmbed{
						{
							Title:       handlers.Cfg.Appearance.BotName + " - Report",
							Description: "Please configure a report channel first. This can be done via the ``/config`` command.",
							Color:       utils.WarnEmbedColour,
						},
					},
				},
			})
			return
		}

		ex = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseModal,
			Data: &discordgo.InteractionResponseData{
				CustomID: "report_create",
				Title:    "Create a report",
				Components: []discordgo.MessageComponent{
					discordgo.ActionsRow{
						Components: []discordgo.MessageComponent{
							discordgo.TextInput{
								CustomID:    "report",
								Label:       "A brief description of your report.",
								Style:       discordgo.TextInputParagraph,
								Placeholder: "Add a brand new meme channel.",
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
		}

		rep_modal_data[i.Member.User.ID] = res["rep_channel"].(string)
	},
}

var reportCreateModal = &handlers.Modal{
	ID: "report_create",
	Exec: func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		desc := i.ModalSubmitData().Components[0].(*discordgo.ActionsRow).Components[0].(*discordgo.TextInput).Value

		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
		})

		// Store the given data
		rep_channel := rep_modal_data[i.Member.User.ID]
		// Remove the data from the map
		delete(rep_modal_data, i.Member.User.ID)

		id := utils.RandomString("r_", 6)

		msg, ex := s.ChannelMessageSendComplex(rep_channel, &discordgo.MessageSend{
			Embeds: []*discordgo.MessageEmbed{
				{
					Author: &discordgo.MessageEmbedAuthor{
						IconURL: i.Member.AvatarURL(""),
						Name:    i.Member.User.Username + "#" + i.Member.User.Discriminator,
					},
					Description: fmt.Sprintf("**Description:** %s\n\n**Status:** Open\n**ID:** %s", desc, id),
					Color:       utils.PlainEmbedColour,
				},
			},
		})

		if ex != nil {
			if strings.Contains(ex.Error(), "Unknown Channel") {
				s.FollowupMessageCreate(s.State.User.ID, i.Interaction, false, &discordgo.WebhookParams{
					Embeds: []*discordgo.MessageEmbed{
						{
							Title:       handlers.Cfg.Appearance.BotName + " - Report",
							Description: "Make sure the configured reports channel still exists and if I have permissions to send message in it. This can be done via the ``/config`` command.",
							Color:       utils.WarnEmbedColour,
						},
					},
				})
			} else {
				utils.Cout("[ERROR] Failed to send message: %v", utils.Red, ex)
				utils.ErrorResponse(s, i.Interaction)
			}
			return
		}

		ex = utils.Firebase.SetFirestore("submissions", id, map[string]interface{}{
			// "guild_id":   i.GuildID,
			"channel_id": msg.ChannelID,
			"message_id": msg.ID,
		}, false)
		if ex != nil {
			utils.Cout("[ERROR] Could not save in Firestore: %v", utils.Red, ex)
			utils.ErrorResponse(s, i.Interaction)
			return
		}

		s.FollowupMessageCreate(s.State.User.ID, i.Interaction, false, &discordgo.WebhookParams{
			Embeds: []*discordgo.MessageEmbed{
				{
					Title:       handlers.Cfg.Appearance.BotName + " - Report",
					Description: "You're report has been submitted!",
					Color:       utils.DefaultEmbedColour,
				},
			},
		})
	},
}
