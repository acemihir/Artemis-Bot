package commands

import "github.com/bwmarrin/discordgo"

func SuggestCommand(s *discordgo.Session, i *discordgo.InteractionCreate) {
	// desc := i.ApplicationCommandData().Options[0].StringValue()
}

func UpvoteButton(s *discordgo.Session, i *discordgo.InteractionCreate) {
	// ...
}

func DownvoteButton(s *discordgo.Session, i *discordgo.InteractionCreate) {
	// ...
}
