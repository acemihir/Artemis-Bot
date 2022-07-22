package commands

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/jerskisnow/Artemis-Bot/src/handlers"
	"github.com/jerskisnow/Artemis-Bot/src/utils"
)

func init() {
	handlers.RegisterCommand(noteCommand)

	handlers.RegisterModal(noteCreateModal)
}

var noteCommand = &handlers.SlashCommand{
	Name: "note",
	Exec: func(s *discordgo.Session, ic *discordgo.InteractionCreate) {
		sbcmd := i.ApplicationCommandData().Options[0].Name
		switch sbcmd {
		case "create":
			commands.noteCreateSubcmd(s, i)
		case "delete":
			commands.noteDeleteSubcmd(s, i)
		case "list":
			commands.noteListSubcmd(s, i)
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
	// ...
}

func noteListSubcmd(s *discordgo.Session, i *discordgo.InteractionCreate) {
	// ...
}

var noteCreateModal = &handlers.Modal{
	ID: "note_create",
	Exec: func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		data := i.ModalSubmitData()
		note := data.Components[0].(*discordgo.ActionsRow).Components[0].(*discordgo.TextInput).Value

		fmt.Println(note)
	},
}
