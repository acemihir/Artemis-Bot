package handlers

import (
	"github.com/bwmarrin/discordgo"
	"github.com/jerskisnow/Suggestions/src/commands"
)

func LinkButton(s *discordgo.Session, i *discordgo.InteractionCreate) {
	data := i.MessageComponentData()

	if data.CustomID == "sug_upvote" {
		commands.UpvoteButton(s, i)
	} else if data.CustomID == "sug_downvote" {
		commands.DownvoteButton(s, i)
	}
}
