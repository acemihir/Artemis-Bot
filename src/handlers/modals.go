package handlers

import (
	"github.com/bwmarrin/discordgo"
)

// ======================
// MODAL HANDLER
// ======================
var mdls = map[string]*Modal{}

type Modal struct {
	ID   string
	Exec func(*discordgo.Session, *discordgo.InteractionCreate)
}

func RegisterModal(modal *Modal) {
	mdls[modal.ID] = modal
}

func LinkModal(s *discordgo.Session, i *discordgo.InteractionCreate) {
	data := i.ModalSubmitData()

	if v, ok := mdls[data.CustomID]; ok {
		v.Exec(s, i)
	}
}
