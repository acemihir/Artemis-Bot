package commands

import "github.com/bwmarrin/discordgo"

func SuggestCommand(s *discordgo.Session, i *discordgo.InteractionCreate) {
	// desc := i.ApplicationCommandData().Options[0].StringValue()

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{},
	})

}

func UpvoteButton(s *discordgo.Session, i *discordgo.InteractionCreate) {
	// ...
}

func DownvoteButton(s *discordgo.Session, i *discordgo.InteractionCreate) {
	// ...
}
