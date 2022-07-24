package commands

import (
	"fmt"

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
							CustomID:    "note",
							Label:       "Your note",
							Style:       discordgo.TextInputParagraph,
							Placeholder: "Buy some strawberries at the store.",
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
		note := data.Components[0].(*discordgo.ActionsRow).Components[0].(*discordgo.TextInput).Value

		fmt.Println(note)
		utils.ComingSoonResponse(s, i.Interaction)
	},
}
