package commands

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/OnlyF0uR/Artemis-Bot/src/handlers"
	"github.com/OnlyF0uR/Artemis-Bot/src/utils"
	"github.com/bwmarrin/discordgo"
)

type modalSuggestionData struct {
	sug_channel    string
	upvote_emoji   string
	downvote_emoji string
}

var sug_modal_data = make(map[string]modalSuggestionData)

func init() {
	handlers.RegisterCommand(suggestCmd)

	handlers.RegisterMessageComponent(upvoteButton)
	handlers.RegisterMessageComponent(downvoteButton)

	handlers.RegisterModal(suggestionCreateModal)
}

var suggestCmd = &handlers.SlashCommand{
	Name: "suggest",
	Exec: func(s *discordgo.Session, i *discordgo.InteractionCreate) {
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
							Title:       handlers.Cfg.Appearance.BotName + " - Suggest",
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
				CustomID: "suggestion_create",
				Title:    "Create a suggestion",
				Components: []discordgo.MessageComponent{
					discordgo.ActionsRow{
						Components: []discordgo.MessageComponent{

							discordgo.TextInput{
								CustomID:    "suggestion",
								Label:       "Your suggestion",
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
			return
		}

		upvote_emoji := "⬆️"
		if res["upvote_emoji"] != nil {
			upvote_emoji = res["upvote_emoji"].(string)
		}
		downvote_emoji := "⬇️"
		if res["downvote_emoji"] != nil {
			downvote_emoji = res["downvote_emoji"].(string)
		}

		sug_modal_data[i.Member.User.ID] = modalSuggestionData{
			sug_channel:    res["sug_channel"].(string),
			upvote_emoji:   upvote_emoji,
			downvote_emoji: downvote_emoji,
		}
	},
}

var suggestionCreateModal = &handlers.Modal{
	ID: "suggestion_create",
	Exec: func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		desc := i.ModalSubmitData().Components[0].(*discordgo.ActionsRow).Components[0].(*discordgo.TextInput).Value

		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
		})

		// Store the given data
		data := sug_modal_data[i.Member.User.ID]
		// Remove the data from the map
		delete(sug_modal_data, i.Member.User.ID)

		id := utils.RandomString("s_", 6)

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
							Style: discordgo.DangerButton,
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
							Title:       handlers.Cfg.Appearance.BotName + " - Suggest",
							Description: "Make sure the configured suggestions channel still exists and if I have permissions to send message in it. This can be done via the ``/config`` command.",
							Color:       utils.WarnEmbedColour,
						},
					},
				})
			} else {
				utils.Cout("[ERROR] Failed to send message: %v", utils.Red, ex)
				utils.ErrorFollowUp(s, i.Interaction)
			}
			return
		}

		ex = utils.Firebase.SetFirestore("submissions", id, map[string]interface{}{
			// "guild_id":   i.GuildID,
			"channel_id": msg.ChannelID,
			"message_id": msg.ID,
			"upvotes":    []string{},
			"downvotes":  []string{},
		}, false)
		if ex != nil {
			utils.Cout("[ERROR] Could not save in Firestore: %v", utils.Red, ex)
			utils.ErrorFollowUp(s, i.Interaction)
			return
		}

		// Add upvotes & downvotes to the cache because most people vote on their own suggestion anyway
		ex = utils.Cache.SetCache(id, "{\"Upvotes\":[],\"Downvotes\":[]}")
		if ex != nil {
			utils.Cout("[ERROR] Could not set in Redis: %v", utils.Red, ex)
			utils.ErrorFollowUp(s, i.Interaction)
			return
		}

		s.FollowupMessageCreate(s.State.User.ID, i.Interaction, false, &discordgo.WebhookParams{
			Embeds: []*discordgo.MessageEmbed{
				{
					Title:       handlers.Cfg.Appearance.BotName + " - Suggest",
					Description: "You're suggestion has been submitted!",
					Color:       utils.DefaultEmbedColour,
				},
			},
		})
	},
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

func cannotVoteTwice(s *discordgo.Session, i *discordgo.InteractionCreate) {
	s.FollowupMessageCreate(s.State.User.ID, i.Interaction, false, &discordgo.WebhookParams{
		Embeds: []*discordgo.MessageEmbed{
			{
				Description: "You cannot vote twice on the same entry.",
				Color:       utils.ErrorEmbedColour,
			},
		},
		Flags: 1 << 6,
	})
}

func countAndMutate(arr *[]interface{}, uid string, rm bool) int {
	var vc int = 0
	for i, v := range *arr {
		if v.(string) == uid {
			if rm {
				tmp := *arr
				*arr = append(tmp[:i], tmp[i+1:]...)
				continue
			} else {
				return -1
			}
		}
		vc++
	}
	if !rm {
		tmp := *arr
		*arr = append(tmp, uid)
		vc++
	}
	return vc
}

var upvoteButton = &handlers.MessageComponent{
	ID: "sug_upvote",
	Exec: func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseDeferredMessageUpdate,
		})

		embed := i.Message.Embeds[0]
		desc_array := strings.Split(embed.Description, "\n")
		id := strings.Split(desc_array[len(desc_array)-3], " ")[1]

		in_cache, ex := utils.Cache.ExistsCache(id)
		if ex != nil {
			utils.Cout("[ERROR] Exists in Redis failed: %v", utils.Red, ex)
			voteError(s, i)
			return
		}

		vote_data := utils.SuggestionVotes{}

		// Check if the data is not in cache
		if in_cache == 0 {
			// Fetch the data from Firestore
			res, ex := utils.Firebase.GetFirestore("submissions", id)
			if ex != nil {
				voteError(s, i)
				utils.Cout("[ERROR] Get from Firestore failed: %v", utils.Red, ex)
				return
			}

			if res["upvotes"] == nil {
				vote_data.Upvotes = []interface{}{}
			} else {
				vote_data.Upvotes = res["upvotes"].([]interface{})
			}
			if res["downvotes"] == nil {
				vote_data.Downvotes = []interface{}{}
			} else {
				vote_data.Downvotes = res["downvotes"].([]interface{})
			}
		} else {
			// Fetch the data from redis
			res, ex := utils.Cache.GetCache(id)
			if ex != nil {
				voteError(s, i)
				utils.Cout("[ERROR] Get from Redis failed: %v", utils.Red, ex)
				return
			}

			// Unmarshal the data into vote_data
			json.Unmarshal([]byte(res), &vote_data)
		}

		upvote_count := countAndMutate(&vote_data.Upvotes, i.Member.User.ID, false)
		if upvote_count == -1 {
			cannotVoteTwice(s, i)
			return
		}
		downvote_count := countAndMutate(&vote_data.Downvotes, i.Member.User.ID, true)

		upvotes_string := "upvotes"
		if upvote_count == 1 {
			upvotes_string = "upvote"
		}
		downvotes_string := "downvotes"
		if downvote_count == 1 {
			downvotes_string = "downvote"
		}

		desc_array[len(desc_array)-1] = fmt.Sprintf("*%d - %s | %d - %s*", upvote_count, upvotes_string, downvote_count, downvotes_string)
		embed.Description = strings.Join(desc_array, "\n")

		res, _ := json.Marshal(vote_data)
		utils.Cache.SetCache(id, string(res))

		s.FollowupMessageEdit(s.State.User.ID, i.Interaction, i.Message.ID, &discordgo.WebhookEdit{
			Embeds: []*discordgo.MessageEmbed{
				embed,
			},
		})
	},
}

var downvoteButton = &handlers.MessageComponent{
	ID: "sug_downvote",
	Exec: func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseDeferredMessageUpdate,
		})

		embed := i.Message.Embeds[0]
		desc_array := strings.Split(embed.Description, "\n")
		id := strings.Split(desc_array[len(desc_array)-3], " ")[1]

		in_cache, ex := utils.Cache.ExistsCache(id)
		if ex != nil {
			voteError(s, i)
			utils.Cout("[ERROR] Exists in Redis failed: %v", utils.Red, ex)
			return
		}

		vote_data := utils.SuggestionVotes{}

		// Check if the data is not in cache
		if in_cache == 0 {
			// Fetch the data from Firestore
			res, ex := utils.Firebase.GetFirestore("submissions", id)
			if ex != nil {
				voteError(s, i)
				utils.Cout("[ERROR] Get from Firestore failed: %v", utils.Red, ex)
				return
			}

			if res["upvotes"] == nil {
				vote_data.Upvotes = []interface{}{}
			} else {
				vote_data.Upvotes = res["upvotes"].([]interface{})
			}
			if res["downvotes"] == nil {
				vote_data.Downvotes = []interface{}{}
			} else {
				vote_data.Downvotes = res["downvotes"].([]interface{})
			}
		} else {
			// Fetch the data from redis
			res, ex := utils.Cache.GetCache(id)
			if ex != nil {
				voteError(s, i)
				utils.Cout("[ERROR] Get from Redis failed: %v", utils.Red, ex)
				return
			}

			// Unmarshal the data into vote_data
			json.Unmarshal([]byte(res), &vote_data)
		}

		downvote_count := countAndMutate(&vote_data.Downvotes, i.Member.User.ID, false)
		if downvote_count == -1 {
			cannotVoteTwice(s, i)
			return
		}
		upvote_count := countAndMutate(&vote_data.Upvotes, i.Member.User.ID, true)

		desc_array[len(desc_array)-1] = fmt.Sprintf("*%d - upvotes | %d - downvotes*", upvote_count, downvote_count)
		embed.Description = strings.Join(desc_array, "\n")

		res, _ := json.Marshal(vote_data)
		utils.Cache.SetCache(id, string(res))

		s.FollowupMessageEdit(s.State.User.ID, i.Interaction, i.Message.ID, &discordgo.WebhookEdit{
			Embeds: []*discordgo.MessageEmbed{
				embed,
			},
		})
	},
}
