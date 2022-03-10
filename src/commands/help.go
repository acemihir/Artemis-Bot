package commands

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/bwmarrin/discordgo"
)

func HelpCommand(s *discordgo.Session, i *discordgo.InteractionCreate) {
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
		log.Fatalln("[ERROR] Failed to fetch application commands. (/help)")
		return
	}

	for _, v := range cmds {
		buffer.WriteString(fmt.Sprintf("/%s", v.Name))

		spacing, ex := strconv.Atoi(os.Getenv("HELP_SPACING_BASE"))
		if ex != nil {
			log.Fatalln("[ERROR] Could not parse HELP_SPACING_BASE. (/help)")
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
}
