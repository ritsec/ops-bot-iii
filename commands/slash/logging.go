package slash

import (
	"github.com/bwmarrin/discordgo"
	"github.com/sirupsen/logrus"
	"gitlab.ritsec.cloud/1nv8rZim/ops-bot-iii/commands/slash/permission"
	"gitlab.ritsec.cloud/1nv8rZim/ops-bot-iii/config"
	"gitlab.ritsec.cloud/1nv8rZim/ops-bot-iii/logging"
	"gitlab.ritsec.cloud/1nv8rZim/ops-bot-iii/structs"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
)

// Log is a slash command that allows users to get or set the logging level
func Log() *structs.SlashCommand {
	return &structs.SlashCommand{
		Command: &discordgo.ApplicationCommand{
			Name:        "log",
			Description: "Get or set the logging level",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "level",
					Description: "Logging Level",
					Required:    false,
					Choices: []*discordgo.ApplicationCommandOptionChoice{
						{
							Name:  "Debug Low",
							Value: logging.LevelNameMap[logging.DebugLowLevel],
						},
						{
							Name:  "Debug",
							Value: logging.LevelNameMap[logging.DebugLevel],
						},
						{
							Name:  "Info",
							Value: logging.LevelNameMap[logging.InfoLevel],
						},
						{
							Name:  "Warn",
							Value: logging.LevelNameMap[logging.WarningLevel],
						},
						{
							Name:  "Error",
							Value: logging.LevelNameMap[logging.ErrorLevel],
						},
						{
							Name:  "Critical",
							Value: logging.LevelNameMap[logging.CriticalLevel],
						},
					},
				},
			},
			DefaultMemberPermissions: &permission.Admin,
		},
		Handler: func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			span := tracer.StartSpan(
				"commands.slash.log:Log",
				tracer.ResourceName("/log"),
			)
			defer span.Finish()

			logging.Debug(s, "Log command received", i.Member.User, span)

			if len(i.ApplicationCommandData().Options) == 0 {
				err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: "Current logging level: " + logging.LevelNameMap[logging.LogLevel()],
						Flags:   discordgo.MessageFlagsEphemeral,
					},
				})
				if err != nil {
					logging.Error(s, "Error sending current log level", i.Member.User, span, logrus.Fields{"error": err})
				}
			} else {
				config.SetLoggingLevel(i.ApplicationCommandData().Options[0].StringValue())

				logging.Critical(s, "Logging level changed to "+i.ApplicationCommandData().Options[0].StringValue(), i.Member.User, span)

				err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: "Logging level changed to " + i.ApplicationCommandData().Options[0].StringValue(),
						Flags:   discordgo.MessageFlagsEphemeral,
					},
				})
				if err != nil {
					logging.Error(s, "Error sending confirmation of change of log level", i.Member.User, span, logrus.Fields{"error": err})
				}
			}
		},
	}
}
