package slash

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/ritsec/ops-bot-iii/commands/slash/permission"
	"github.com/ritsec/ops-bot-iii/logging"
	"github.com/ritsec/ops-bot-iii/structs"
	"github.com/sirupsen/logrus"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
)

// Purge is the purge command
func Purge() *structs.SlashCommand {
	min := float64(1)
	return &structs.SlashCommand{
		Command: &discordgo.ApplicationCommand{
			Name:        "purge",
			Description: "Purge messages from a channel",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionInteger,
					Name:        "amount",
					Description: "The amount of messages to purge (default: 100)",
					Required:    false,
					MinValue:    &min,
					MaxValue:    100,
				},
			},
			DefaultMemberPermissions: &permission.Admin,
		},
		Handler: func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			span := tracer.StartSpan(
				"commands.slash.purge:Purge",
				tracer.ResourceName("/purge"),
			)
			defer span.Finish()

			logging.Debug(s, "Purge command received", i.Member.User, span)

			var (
				message_ids []string
				messages    int64
			)

			if len(i.ApplicationCommandData().Options) != 0 {
				messages = i.ApplicationCommandData().Options[0].IntValue()
			} else {
				messages = 100
			}

			raw_messages, err := s.ChannelMessages(i.ChannelID, int(messages), "", "", "")
			if err != nil {
				logging.Error(s, err.Error(), i.Member.User, span, logrus.Fields{"error": err})
			}

			for _, message := range raw_messages {
				message_ids = append(message_ids, message.ID)
			}

			err = s.ChannelMessagesBulkDelete(i.ChannelID, message_ids)
			if err != nil {
				logging.Error(s, err.Error(), i.Member.User, span, logrus.Fields{"error": err})
			}

			err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "Purged " + fmt.Sprint(len(raw_messages)) + " messages!",
					Flags:   discordgo.MessageFlagsEphemeral,
				},
			})

			logging.Debug(s, "Purged "+fmt.Sprint(len(raw_messages))+" messages!", i.Member.User, span)
			if err != nil {
				logging.Error(s, err.Error(), i.Member.User, span, logrus.Fields{"error": err})
			}
		},
	}
}
