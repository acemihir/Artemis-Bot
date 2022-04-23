package handlers

import (
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
	case "status_change":
		commands.StatusDropdown(s, i)

	// Config part
	case "cfg_main_auth":
		commands.ConfigMainAuthButton(s, i)
	case "cfg_main_chns":
		commands.ConfigMainChannelsButton(s, i)
	case "cfg_main_appear":
		commands.ConfigMainAppearanceButton(s, i)
		// case "cfg_main_misc":
		// 	commands.ConfigMainMiscButton(s, i)
	case "cfg_auth_staffrole":
		commands.ConfigAuthStaffroleButton(s, i)
	case "cfg_chns_sug":
		commands.ConfigChnsSugButton(s, i)
	case "cfg_chns_rep":
		commands.ConfigChnsRepButton(s, i)
	case "cfg_appear_upvote":
		commands.ConfigAppearUpvoteButton(s, i)
	case "cfg_appear_downvote":
		commands.ConfigAppearDownvoteButton(s, i)
	}
}
