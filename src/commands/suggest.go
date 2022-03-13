package commands

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/jerskisnow/Artemis-Bot/src/utils"
)

func SuggestCommand(s *discordgo.Session, i *discordgo.InteractionCreate) {
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
	})

	res, ex := utils.Firebase.GetFirestore("guilds", i.GuildID)
	if ex != nil {
		utils.Cout("[ERROR] Get from Firestore failed: %v", utils.Red, ex)
		utils.ErrorResponse(s, i.Interaction)
		return
	}

	if len(res) == 0 || res["sug_channel"] == nil {
		s.FollowupMessageCreate(s.State.User.ID, i.Interaction, false, &discordgo.WebhookParams{
			Embeds: []*discordgo.MessageEmbed{
				{
					Title:       "Artemis - Suggest",
					Description: "Please configure a suggestion channel first. This can be done via the ``/config`` command.",
					Color:       utils.WarnEmbedColour,
				},
			},
		})
		return
	}

	desc := i.ApplicationCommandData().Options[0].StringValue()
	id := utils.CreateId("s_", 6)

	upvote_emoji := "⬆️"
	if res["upvote_emoji"] != nil {
		upvote_emoji = res["upvote_emoji"].(string)
	}
	downvote_emoji := "⬇️"
	if res["downvote_emoji"] != nil {
		downvote_emoji = res["downvote_emoji"].(string)
	}

	msg, ex := s.ChannelMessageSendComplex(res["sug_channel"].(string), &discordgo.MessageSend{
		Embeds: []*discordgo.MessageEmbed{
			{
				Author: &discordgo.MessageEmbedAuthor{
					IconURL: i.Member.AvatarURL(""),
					Name:    i.Member.User.Username + "#" + i.Member.User.Discriminator,
				},
				Description: fmt.Sprintf("**Description:** %s\n\n**Status:** Open\n**ID:** %s\n\n*0 - upvotes | 0 - downvotes*", desc, id),
				Color:       utils.PlainEmbedColour,
			},
		},
		Components: []discordgo.MessageComponent{
			discordgo.ActionsRow{
				Components: []discordgo.MessageComponent{
					discordgo.Button{
						CustomID: "sug_upvote",
						Emoji: discordgo.ComponentEmoji{
							Name: upvote_emoji,
						},
						Label: "Upvote",
						Style: discordgo.SuccessButton,
					},
					discordgo.Button{
						CustomID: "sug_downvote",
						Emoji: discordgo.ComponentEmoji{
							Name: downvote_emoji,
						},
						Label: "Downvote",
						Style: discordgo.SuccessButton,
					},
				},
			},
		},
	})

	if ex != nil {
		if strings.Contains(ex.Error(), "Unknown Channel") {
			s.FollowupMessageCreate(s.State.User.ID, i.Interaction, false, &discordgo.WebhookParams{
				Embeds: []*discordgo.MessageEmbed{
					{
						Title:       "Artemis - Suggest",
						Description: "Make sure the configured suggestions channel still exists and if I have permissions to send message in it. This can be done via the ``/config`` command.",
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
		"upvotes":    0,
		"downvotes":  0,
	})
	if ex != nil {
		utils.Cout("[ERROR] Could not save in Firestore: %v", utils.Red, ex)
		utils.ErrorResponse(s, i.Interaction)
		return
	}

	s.FollowupMessageCreate(s.State.User.ID, i.Interaction, false, &discordgo.WebhookParams{
		Embeds: []*discordgo.MessageEmbed{
			{
				Title:       "Artemis - Suggest",
				Description: "You're suggestion has been submitted!",
				Color:       utils.DefaultEmbedColour,
			},
		},
	})
}

func UpvoteButton(s *discordgo.Session, i *discordgo.InteractionCreate) {
	embed := i.Message.Embeds[0]
	desc_array := strings.Split(embed.Description, "\n")
	id := strings.Split(desc_array[len(desc_array)-3], " ")[1]

	b, ex := utils.Cache.ExistsCache(id)
	if ex != nil {
		// TODO
	}

	if b == 0 {
		// Set in cache
		ex = utils.Cache.SetCache(id, "1:0")
		if ex != nil {
			// TODO
		}
	} else {
		// Get from cache
		res, ex := utils.Cache.GetCache(id)
		if ex != nil {
			// TODO
		}

		vote_array := strings.Split(res, ":")

		vote_n, ex := strconv.Atoi(vote_array[0])
		if ex != nil {
			// TODO
		}

		vote_n++
		ex = utils.Cache.SetCache(id, strconv.Itoa(vote_n)+":"+vote_array[1])
		if ex != nil {
			// TODO
		}
	}

	// TODO: Update the mesasge

	// desc_array[len(desc_array) - 1] = "*0 - upvotes | 0 - downvotes*"

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseUpdateMessage,
	})
}

func DownvoteButton(s *discordgo.Session, i *discordgo.InteractionCreate) {
	embed := i.Message.Embeds[0]
	desc_array := strings.Split(embed.Description, "\n")
	id := strings.Split(desc_array[len(desc_array)-3], " ")[1]

	b, ex := utils.Cache.ExistsCache(id)
	if ex != nil {
		// TODO
	}

	if b == 0 {
		// Set in cache
		ex = utils.Cache.SetCache(id, "0:1")
		if ex != nil {
			// TODO
		}
	} else {
		// Get from cache
		res, ex := utils.Cache.GetCache(id)
		if ex != nil {
			// TODO
		}

		vote_array := strings.Split(res, ":")

		vote_n, ex := strconv.Atoi(vote_array[1])
		if ex != nil {
			// TODO
		}

		vote_n++
		ex = utils.Cache.SetCache(id, vote_array[0]+":"+strconv.Itoa(vote_n))
		if ex != nil {
			// TODO
		}
	}

	// TODO: Update the mesasge

	// desc_array[len(desc_array) - 1] = "*0 - upvotes | 0 - downvotes*"

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseUpdateMessage,
	})
}
