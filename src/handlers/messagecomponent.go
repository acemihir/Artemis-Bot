package handlers

import (
	"github.com/bwmarrin/discordgo"
)

// =========================
// MESSAGE COMPONENT HANDLER
// =========================
var msgCmpnts = map[string]*MessageComponent{}

type MessageComponent struct {
	ID   string
	Exec func(*discordgo.Session, *discordgo.InteractionCreate)
}

func RegisterMessageComponent(cmp *MessageComponent) {
	msgCmpnts[cmp.ID] = cmp
}

func LinkMessageComponent(s *discordgo.Session, i *discordgo.InteractionCreate) {
	data := i.MessageComponentData()

	if v, ok := mdls[data.CustomID]; ok {
		v.Exec(s, i)
	}
}
