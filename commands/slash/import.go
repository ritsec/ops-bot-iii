package slash

import (
	"github.com/bwmarrin/discordgo"
	"github.com/ritsec/ops-bot-iii/commands/slash/permission"
	"github.com/ritsec/ops-bot-iii/logging"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
)

// Import slash command
func Import() (*discordgo.ApplicationCommand, func(s *discordgo.Session, i *discordgo.InteractionCreate)) {
	return &discordgo.ApplicationCommand{
			Name:                     "import",
			Description:              "Import signins from csv",
			DefaultMemberPermissions: &permission.Admin,
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionAttachment,
					Name:        "file",
					Description: "A csv file containing sign-ins to import.",
					Required:    true,
				},
			},
		},
		func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			span := tracer.StartSpan(
				"commands.slash.import:Import",
				tracer.ResourceName("/import"),
			)
			defer span.Finish()

			logging.Debug(s, "Import command received", i.Member.User, span)

			err := s.InteractionRespond(
				i.Interaction,
				&discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: importMessage(i.Member.User.ID, span.Context()),
						Flags:   discordgo.MessageFlagsEphemeral,
					},
				},
			)
			if err != nil {
				logging.Error(s, err.Error(), i.Member.User, span)
			}

			data := i.ApplicationCommandData()
			options := data.Options
		}
}
