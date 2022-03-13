package handlers

import (
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/jerskisnow/Artemis-Bot/src/commands"
)

func LinkButton(s *discordgo.Session, i *discordgo.InteractionCreate) {
	data := i.MessageComponentData()

	switch data.CustomID {
	case "sug_upvote":
		commands.UpvoteButton(s, i)
	case "sug_downvote":
		commands.DownvoteButton(s, i)
	default:
		if strings.HasPrefix(data.CustomID, "status-") {
			commands.StatusDropdown(s, i)
		}
	}
}
