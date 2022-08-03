package commands

import (
	"encoding/base64"

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

		// Save the note in firestore
		ex = utils.Firebase.SetFirestore("notes", title, map[string]interface{}{
			"author":   i.User.ID,
			"contents": ct,
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
					Title:       "Artemis - Suggest",
					Description: "Note has been submitted.",
					Color:       utils.DefaultEmbedColour,
				},
			},
		})
	},
}
