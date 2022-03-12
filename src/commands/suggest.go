package commands

import (
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/jerskisnow/Artemis-Bot/src/utils"
)

func SuggestCommand(s *discordgo.Session, i *discordgo.InteractionCreate) {
	// desc := i.ApplicationCommandData().Options[0].StringValue()

	// Defer the message while we handle everything
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{},
	})

	// Fetch data from firestorm
	res, ex := utils.Firebase.GetFirestore("guilds", i.GuildID)
	if ex != nil {
		utils.Cout("[ERROR] Get from Firestore failed: %v", utils.Red, ex)
		utils.ErrorResponse(s, i.Interaction)
		return
	}

	if len(res) == 0 || res["sug_channel"] == nil {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Embeds: []*discordgo.MessageEmbed{
					{
						Title:       "Artemis - Suggest",
						Description: "Please configure a suggestion channel first. This can be done via the ``/config`` command.",
						Color:       0xffcb47,
					},
				},
			},
		})
		return
	}

	// c, ex := s.Channel(res["sug_channel"])
	// if ex != nil {
	// 	utils.Cout("[ERROR] Failed to get channel: %v", utils.Red, ex)
	// 	utils.ErrorResponse(s, i.Interaction)
	// }

	// TODO: Fix description
	msg, ex := s.ChannelMessageSendComplex(res["sug_channel"].(string), &discordgo.MessageSend{
		Embeds: []*discordgo.MessageEmbed{
			{
				Author: &discordgo.MessageEmbedAuthor{
					IconURL: i.User.AvatarURL(""),
					Name:    i.User.Username + "#" + i.User.Discriminator,
				},
				Description: "TODO",
				Color:       0x614832,
			},
		},
	})

	if ex != nil {
		if strings.Contains(ex.Error(), "Unknown Channel") {
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Embeds: []*discordgo.MessageEmbed{
						{
							Title:       "Artemis - Suggest",
							Description: "Make sure the configured suggestions channel still exists and if I have permissions to send message in it. This can be done via the ``/config`` command.",
							Color:       0xffcb47,
						},
					},
				},
			})
		} else {
			utils.Cout("[ERROR] Failed to send message: %v", utils.Red, ex)
			utils.ErrorResponse(s, i.Interaction)
		}
		return
	}

	// TODO Send message that the suggestion has been submitted
}

func UpvoteButton(s *discordgo.Session, i *discordgo.InteractionCreate) {
	// ...
}

func DownvoteButton(s *discordgo.Session, i *discordgo.InteractionCreate) {
	// ...
}
