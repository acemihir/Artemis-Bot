package utils

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

const (
	Reset  = "\033[0m"
	Red    = "\033[31m"
	Green  = "\033[32m"
	Yellow = "\033[33m"
	Blue   = "\033[34m"
	Purple = "\033[35m"
	Cyan   = "\033[36m"
	Gray   = "\033[37m"
	White  = "\033[97m"

	ErrorEmbedColour   = 0xff4a4a
	WarnEmbedColour    = 0xffcb47
	DefaultEmbedColour = 0x614832
	PlainEmbedColour   = 0x2f3136

	GreenEmbedColour  = 0x97ff78
	YellowEmbedColour = 0xffed78
	RedEmbedColour    = 0xfc5d5d

	AdminPermission = 0x8    // ADMINISTRATOR
	StaffPermission = 0x2000 // MANAGE_MESSAGES
)

type SuggestionVotes struct {
	Upvotes   []interface{}
	Downvotes []interface{}
}

func Cout(text string, colour string, params ...interface{}) {
	if len(params) == 0 {
		fmt.Println(colour + text + Reset)
	} else {
		fmt.Printf(colour+text+Reset+"\n", params...)
	}
}

func ErrorResponse(s *discordgo.Session, i *discordgo.Interaction) {
	s.InteractionRespond(i, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{
				{
					Description: "Oops! A wild error seems to have occured.\n\nPlease try again later, if this error is persistent please report it in our Support discord.",
					Color:       ErrorEmbedColour,
				},
			},
			Components: []discordgo.MessageComponent{
				discordgo.ActionsRow{
					Components: []discordgo.MessageComponent{
						discordgo.Button{
							Emoji: discordgo.ComponentEmoji{
								Name: "👥",
							},
							Label: "Support",
							Style: discordgo.LinkButton,
							URL:   "https://discord.gg/3SYg3M5",
						},
					},
				},
			},
		},
	})
}

func ComingSoonResponse(s *discordgo.Session, i *discordgo.Interaction) {
	s.InteractionRespond(i, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{
				{
					Description: "This feature will be available soon. Feel free to join our support server for more information.",
					Color:       DefaultEmbedColour,
				},
			},
			Components: []discordgo.MessageComponent{
				discordgo.ActionsRow{
					Components: []discordgo.MessageComponent{
						discordgo.Button{
							Emoji: discordgo.ComponentEmoji{
								Name: "👥",
							},
							Label: "Info",
							Style: discordgo.LinkButton,
							URL:   "https://discord.gg/3SYg3M5",
						},
					},
				},
			},
		},
	})
}

func HasPermission(perms int64, req int64) bool {
	return (perms & req) == req
}
