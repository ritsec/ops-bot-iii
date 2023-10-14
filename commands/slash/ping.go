package slash

import (
	"github.com/bwmarrin/discordgo"
	"github.com/ritsec/ops-bot-iii/commands/slash/permission"
	"github.com/ritsec/ops-bot-iii/logging"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
)

// Ping is a slash command that responds with "Pong!"
func Ping() (*discordgo.ApplicationCommand, func(s *discordgo.Session, i *discordgo.InteractionCreate)) {
	return &discordgo.ApplicationCommand{
			Name:                     "ping",
			Description:              "Pong!",
			DefaultMemberPermissions: &permission.Member,
		},
		func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			span := tracer.StartSpan(
				"commands.slash.ping:Ping",
				tracer.ResourceName("/ping"),
			)
			defer span.Finish()

			logging.Debug(s, "Ping command received", i.Member.User, span)

			err := s.InteractionRespond(
				i.Interaction,
				&discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: "Pong!",
						Flags:   discordgo.MessageFlagsEphemeral,
					},
				},
			)
			if err != nil {
				logging.Error(s, err.Error(), i.Member.User, span)
			} else {
				logging.Debug(s, "Pong!", i.Member.User, span)
			}
		}
}
