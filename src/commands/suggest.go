package commands

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/jerskisnow/Artemis-Bot/src/utils"
)

type ModalSuggestionData struct {
	sug_channel    string
	upvote_emoji   string
	downvote_emoji string
}

var modal_data = make(map[string]ModalSuggestionData)

func SuggestCommand(s *discordgo.Session, i *discordgo.InteractionCreate) {
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
			CustomID: "modals_suggestion",
			Title:    "Suggestion - Create",
			Components: []discordgo.MessageComponent{
				discordgo.ActionsRow{
					Components: []discordgo.MessageComponent{
						discordgo.TextInput{
							CustomID:    "suggestion",
							Label:       "A brief description of your suggestion.",
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

	upvote_emoji := "⬆️"
	if res["upvote_emoji"] != nil {
		upvote_emoji = res["upvote_emoji"].(string)
	}
	downvote_emoji := "⬇️"
	if res["downvote_emoji"] != nil {
		downvote_emoji = res["downvote_emoji"].(string)
	}

	modal_data[i.Member.User.ID] = ModalSuggestionData{
		sug_channel:    res["sug_channel"].(string),
		upvote_emoji:   upvote_emoji,
		downvote_emoji: downvote_emoji,
	}
}

func SuggestionModal(s *discordgo.Session, i *discordgo.InteractionCreate) {
	desc := i.ModalSubmitData().Components[0].(*discordgo.ActionsRow).Components[0].(*discordgo.TextInput).Value

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
	})

	// Store the given data
	data := modal_data[i.Member.User.ID]
	// Remove the data from the map
	delete(modal_data, i.Member.User.ID)

	id := utils.CreateId("s_", 6)

	msg, ex := s.ChannelMessageSendComplex(data.sug_channel, &discordgo.MessageSend{
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
							Name: data.upvote_emoji,
						},
						Label: "Upvote",
						Style: discordgo.SuccessButton,
					},
					discordgo.Button{
						CustomID: "sug_downvote",
						Emoji: discordgo.ComponentEmoji{
							Name: data.downvote_emoji,
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

func voteError(s *discordgo.Session, i *discordgo.InteractionCreate) {
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
}

// ===========================================
// TODO: Check if the user already voted
func UpvoteButton(s *discordgo.Session, i *discordgo.InteractionCreate) {
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseDeferredMessageUpdate,
	})

	embed := i.Message.Embeds[0]
	desc_array := strings.Split(embed.Description, "\n")
	id := strings.Split(desc_array[len(desc_array)-3], " ")[1]

	b, ex := utils.Cache.ExistsCache(id)
	if ex != nil {
		voteError(s, i)
		utils.Cout("[ERROR] Exists in redis failed: %v", utils.Red, ex)
		return
	}

	var upvotes string = "0"
	var downvotes string = "0"

	if b == 0 {
		ex = utils.Cache.SetCache(id, "1:0")
		if ex != nil {
			voteError(s, i)
			utils.Cout("[ERROR] Set in redis failed: %v", utils.Red, ex)
			return
		}
		upvotes = "1"
	} else {
		// Get from cache
		res, ex := utils.Cache.GetCache(id)
		if ex != nil {
			voteError(s, i)
			utils.Cout("[ERROR] Get from redis failed: %v", utils.Red, ex)
			return
		}

		// Split the voting data into an array
		vote_array := strings.Split(res, ":")
		// Parse the voting data to an integer
		vote_n, ex := strconv.Atoi(vote_array[0])
		if ex != nil {
			voteError(s, i)
			utils.Cout("[ERROR] Parse votecount to int failed: %v", utils.Red, ex)
			return
		}

		// Incremeant the voting data
		vote_n++

		upvotes = strconv.Itoa(vote_n)
		downvotes = vote_array[1]

		ex = utils.Cache.SetCache(id, upvotes+":"+downvotes)
		if ex != nil {
			voteError(s, i)
			utils.Cout("[ERROR] Set in redis failed: %v", utils.Red, ex)
			return
		}
	}

	desc_array[len(desc_array)-1] = fmt.Sprintf("*%s - upvotes | %s - downvotes*", upvotes, downvotes)
	embed.Description = strings.Join(desc_array, "\n")

	s.FollowupMessageEdit(s.State.User.ID, i.Interaction, i.Message.ID, &discordgo.WebhookEdit{
		Embeds: []*discordgo.MessageEmbed{
			embed,
		},
	})
}

// TODO: Check if the user already voted
func DownvoteButton(s *discordgo.Session, i *discordgo.InteractionCreate) {
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseDeferredMessageUpdate,
	})

	embed := i.Message.Embeds[0]
	desc_array := strings.Split(embed.Description, "\n")
	id := strings.Split(desc_array[len(desc_array)-3], " ")[1]

	b, ex := utils.Cache.ExistsCache(id)
	if ex != nil {
		voteError(s, i)
		utils.Cout("[ERROR] Exists in redis failed: %v", utils.Red, ex)
		return
	}

	var upvotes string = "0"
	var downvotes string = "0"

	if b == 0 {
		ex = utils.Cache.SetCache(id, "0:1")
		if ex != nil {
			voteError(s, i)
			utils.Cout("[ERROR] Set in redis failed: %v", utils.Red, ex)
			return
		}
		downvotes = "1"
	} else {
		// Get from cache
		res, ex := utils.Cache.GetCache(id)
		if ex != nil {
			voteError(s, i)
			utils.Cout("[ERROR] Get from redis failed: %v", utils.Red, ex)
			return
		}

		// Split the voting data into an array
		vote_array := strings.Split(res, ":")
		// Parse the voting data to an integer
		vote_n, ex := strconv.Atoi(vote_array[1])
		if ex != nil {
			voteError(s, i)
			utils.Cout("[ERROR] Parse votecount to int failed: %v", utils.Red, ex)
			return
		}

		// Incremeant the voting data
		vote_n++

		upvotes = vote_array[0]
		downvotes = strconv.Itoa(vote_n)

		ex = utils.Cache.SetCache(id, upvotes+":"+downvotes)
		if ex != nil {
			voteError(s, i)
			utils.Cout("[ERROR] Set in redis failed: %v", utils.Red, ex)
			return
		}
	}

	desc_array[len(desc_array)-1] = fmt.Sprintf("*%s - upvotes | %s - downvotes*", upvotes, downvotes)
	embed.Description = strings.Join(desc_array, "\n")

	s.FollowupMessageEdit(s.State.User.ID, i.Interaction, i.Message.ID, &discordgo.WebhookEdit{
		Embeds: []*discordgo.MessageEmbed{
			embed,
		},
	})
}
