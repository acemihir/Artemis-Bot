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
)

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
