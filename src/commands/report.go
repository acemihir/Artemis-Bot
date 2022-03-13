package commands

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/jerskisnow/Artemis-Bot/src/utils"
)

func ReportCommand(s *discordgo.Session, i *discordgo.InteractionCreate) {
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{},
	})

	res, ex := utils.Firebase.GetFirestore("guilds", i.GuildID)
	if ex != nil {
		utils.Cout("[ERROR] Get from Firestore failed: %v", utils.Red, ex)
		utils.ErrorResponse(s, i.Interaction)
		return
	}

	if len(res) == 0 || res["rep_channel"] == nil {
		s.FollowupMessageCreate(s.State.User.ID, i.Interaction, false, &discordgo.WebhookParams{
			Embeds: []*discordgo.MessageEmbed{
				{
					Title:       "Artemis - Report",
					Description: "Please configure a report channel first. This can be done via the ``/config`` command.",
					Color:       utils.WarnEmbedColour,
				},
			},
		})
		return
	}

	desc := i.ApplicationCommandData().Options[0].StringValue()
	id := utils.CreateId("r_", 6)

	msg, ex := s.ChannelMessageSendComplex(res["rep_channel"].(string), &discordgo.MessageSend{
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
						Title:       "Artemis - Report",
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
	})
	if ex != nil {
		utils.Cout("[ERROR] Could not save in Firestore: %v", utils.Red, ex)
		utils.ErrorResponse(s, i.Interaction)
		return
	}

	s.FollowupMessageCreate(s.State.User.ID, i.Interaction, false, &discordgo.WebhookParams{
		Embeds: []*discordgo.MessageEmbed{
			{
				Title:       "Artemis - Report",
				Description: "You're report has been submitted!",
				Color:       utils.DefaultEmbedColour,
			},
		},
	})
}
