package handlers

import (
	"github.com/bwmarrin/discordgo"
)

func LinkButton(s *discordgo.Session, i *discordgo.InteractionCreate) {
	data := i.MessageComponentData()

	if data.CustomID == "sug_upvote" {
		// ...
	} else if data.CustomID == "sug_downvote" {
		// ...
	}
}
