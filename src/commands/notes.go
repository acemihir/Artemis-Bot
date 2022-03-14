package commands

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/jerskisnow/Artemis-Bot/src/utils"
)

func NotesCreateCommand(s *discordgo.Session, i *discordgo.InteractionCreate) {
	ex := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseModal,
		Data: &discordgo.InteractionResponseData{
			CustomID: "modals_notes",
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

func NotesDeleteCommand(s *discordgo.Session, i *discordgo.InteractionCreate) {

}

func NotesListCommand(s *discordgo.Session, i *discordgo.InteractionCreate) {

}

func NotesModal(s *discordgo.Session, i *discordgo.InteractionCreate) {
	data := i.ModalSubmitData()
	note := data.Components[0].(*discordgo.ActionsRow).Components[0].(*discordgo.TextInput).Value

	fmt.Println(note)
}
