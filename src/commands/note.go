package commands

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"time"

	"github.com/OnlyF0uR/Artemis-Bot/src/handlers"
	"github.com/OnlyF0uR/Artemis-Bot/src/utils"
	"github.com/bwmarrin/discordgo"
)

func init() {
	handlers.RegisterCommand(noteCommand)

	handlers.RegisterModal(noteCreateModal)
}

var noteCommand = &handlers.SlashCommand{
	Name: "note",
	Exec: func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		sbcmd := i.ApplicationCommandData().Options[0].Name
		switch sbcmd {
		case "create":
			noteCreateSubcmd(s, i)
		case "delete":
			noteDeleteSubcmd(s, i)
		case "list":
			noteListSubcmd(s, i)
		}
	},
}

func noteCreateSubcmd(s *discordgo.Session, i *discordgo.InteractionCreate) {
	ex := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseModal,
		Data: &discordgo.InteractionResponseData{
			CustomID: "note_create",
			Title:    "Create a note",
			Components: []discordgo.MessageComponent{
				discordgo.ActionsRow{
					Components: []discordgo.MessageComponent{
						discordgo.TextInput{
							CustomID:    "title",
							Label:       "Title",
							Style:       discordgo.TextInputShort,
							Placeholder: "Note #1",
							Required:    true,
							MaxLength:   30,
							MinLength:   1,
						},
						discordgo.TextInput{
							CustomID:    "note",
							Label:       "Your note",
							Style:       discordgo.TextInputParagraph,
							Placeholder: "Buy some strawberries at the store.",
							Required:    true,
							MaxLength:   600,
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

func noteDeleteSubcmd(s *discordgo.Session, i *discordgo.InteractionCreate) {
	note_id := i.ApplicationCommandData().Options[0].StringValue()

	res, ex := utils.Firebase.GetFirestore("notes", note_id)
	if ex != nil {
		utils.Cout("[ERROR] Get from Firestore failed: %v", utils.Red, ex)
		utils.ErrorResponse(s, i.Interaction)
		return
	}

	if len(res) == 0 {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Embeds: []*discordgo.MessageEmbed{
					{
						Title:       handlers.Cfg.Appearance.BotName + " - Notes",
						Description: "You do not have any notes with that title.",
						Color:       utils.WarnEmbedColour,
					},
				},
			},
		})
		return
	}

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{
				{
					Description: "Note deleted.",
					Color:       utils.DefaultEmbedColour,
				},
			},
		},
	})

	utils.ComingSoonResponse(s, i.Interaction)
}

func noteListSubcmd(s *discordgo.Session, i *discordgo.InteractionCreate) {
	utils.ComingSoonResponse(s, i.Interaction)
}

var noteCreateModal = &handlers.Modal{
	ID: "note_create",
	Exec: func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		data := i.ModalSubmitData()

		title := data.Components[0].(*discordgo.ActionsRow).Components[0].(*discordgo.TextInput).Value
		contents := data.Components[0].(*discordgo.ActionsRow).Components[1].(*discordgo.TextInput).Value

		// Defer response
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
		})

		// Let's atleast make sure Google cannot read the notes :/
		key, ex := base64.URLEncoding.DecodeString(handlers.Cfg.Data.EncryptionKey)
		if ex != nil {
			utils.Cout("[ERROR] Could decode base64 encryption key: %v", utils.Red, ex)
			utils.ErrorFollowUp(s, i.Interaction)
			return
		}
		ct := utils.EncryptAES(key, contents)

		// Create hash of author
		h := hmac.New(sha256.New, []byte(handlers.Cfg.Data.HMACKey))
		h.Write([]byte(i.User.ID))

		author_hash := hex.EncodeToString(h.Sum(nil))

		// Generate note id
		note_id := utils.RandomString("n_", 6)

		// Save the note in firestore
		ex = utils.Firebase.SetFirestore("notes", note_id, map[string]interface{}{
			"author_hash": author_hash,
			"title":       title,
			"contents":    ct,
			"timestamp":   time.Now().Unix(),
		}, false)
		if ex != nil {
			utils.Cout("[ERROR] Could not save in Firestore: %v", utils.Red, ex)
			utils.ErrorFollowUp(s, i.Interaction)
			return
		}

		// Notify user about submitted note
		s.FollowupMessageCreate(s.State.User.ID, i.Interaction, false, &discordgo.WebhookParams{
			Embeds: []*discordgo.MessageEmbed{
				{
					Title:       handlers.Cfg.Appearance.BotName + " - Notes",
					Description: "Note has been submitted.",
					Color:       utils.DefaultEmbedColour,
				},
			},
		})
	},
}
