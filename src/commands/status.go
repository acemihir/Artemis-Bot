package commands

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/jerskisnow/Artemis-Bot/src/utils"
)

func StatusCommand(s *discordgo.Session, i *discordgo.InteractionCreate) {
	id := i.ApplicationCommandData().Options[0].StringValue()

	var ex error
	if strings.HasPrefix(id, "s_") {
		ex = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Components: []discordgo.MessageComponent{
					discordgo.ActionsRow{
						Components: []discordgo.MessageComponent{
							discordgo.SelectMenu{
								CustomID:    "status-" + id,
								Placeholder: "Select a status",
								Options: []discordgo.SelectMenuOption{
									{
										Label: "Open",
										Value: "open",
										Emoji: discordgo.ComponentEmoji{
											Name: "⚪",
										},
										Default: true,
									},
									{
										Label: "Considering",
										Value: "considering",
										Emoji: discordgo.ComponentEmoji{
											Name: "🟡",
										},
										Default: false,
									},
									{
										Label: "Approve",
										Value: "approve",
										Emoji: discordgo.ComponentEmoji{
											Name: "🟢",
										},
										Default: false,
									},
									{
										Label: "Reject",
										Value: "reject",
										Emoji: discordgo.ComponentEmoji{
											Name: "🔴",
										},
										Default: false,
									},
								},
							},
						},
					},
				},
			},
		})
	} else if strings.HasPrefix(id, "r_") {
		// ...
	}

	if ex != nil {
		utils.Cout("[ERROR] Could not open up the modal: %v", utils.Red, ex)
		utils.ErrorResponse(s, i.Interaction)
		return
	}
}

func StatusDropdown(s *discordgo.Session, i *discordgo.InteractionCreate) {
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseDeferredMessageUpdate,
	})

	data := i.MessageComponentData()

	id := strings.Split(data.CustomID, "-")[1]
	status := data.Values[0]

	fmt.Println(id)
	fmt.Println(status)

	// ...
}
