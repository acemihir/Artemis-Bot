package commands

import (
	"bytes"
	"fmt"
	"os"
	"strconv"

	"github.com/OnlyF0uR/Artemis-Bot/src/handlers"
	"github.com/OnlyF0uR/Artemis-Bot/src/utils"
	"github.com/bwmarrin/discordgo"
)

func init() {
	handlers.RegisterCommand(helpCmd)
}

var helpCmd = &handlers.SlashCommand{
	Name: "help",
	Exec: func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		var buffer bytes.Buffer
		buffer.WriteString("```asciidoc\n")
		buffer.WriteString("== Help == \n[View the autocompletion for more detailed explanation.]\n\n== Commands ==\n")

		var cmds []*discordgo.ApplicationCommand
		var ex error

		if os.Getenv("PRODUCTION") == "0" {
			cmds, ex = s.ApplicationCommands(s.State.User.ID, os.Getenv("GUILD_ID"))
		} else {
			cmds, ex = s.ApplicationCommands(s.State.User.ID, "")
		}

		if ex != nil {
			utils.Cout("[ERROR][CMD-HELP] Failed to fetch application commands: %v", utils.Red, ex)
			return
		}

		for _, v := range cmds {
			buffer.WriteString(fmt.Sprintf("/%s", v.Name))

			spacing, ex := strconv.Atoi(os.Getenv("HELP_SPACING_BASE"))
			if ex != nil {
				utils.Cout("[ERROR][CMD-HELP] Could not parse HELP_SPACING_BASE: %v", utils.Red, ex)
			}
			spacing -= len(v.Name)

			for i := 1; i <= spacing; i++ {
				buffer.WriteString(" ")
			}

			buffer.WriteString(fmt.Sprintf(":: %s\n", v.Description))
		}

		buffer.WriteString("```")

		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: buffer.String(),
			},
		})
	},
}
