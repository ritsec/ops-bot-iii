package slash

import (
	"fmt"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/ritsec/ops-bot-iii/commands/slash/permission"
	"github.com/ritsec/ops-bot-iii/logging"
	"github.com/sirupsen/logrus"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
)

// Purge is the purge command
func Purge() (*discordgo.ApplicationCommand, func(s *discordgo.Session, i *discordgo.InteractionCreate)) {
	min := float64(1)
	return &discordgo.ApplicationCommand{
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
		func(s *discordgo.Session, i *discordgo.InteractionCreate) {
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

			file := fmt.Sprintf("Record of the purge on %v", time.Now())
			file += "-------------------------------"

			for _, message := range raw_messages {
				message_ids = append(message_ids, message.ID)
				// Timestamp "may be" removed in a future API version. Too bad!
				file += fmt.Sprintf("\n%v SENT AT %v (EDITED AT %v)", message.Author, message.Timestamp, message.EditedTimestamp)
				file += fmt.Sprintf("%v", message.Content)
			}

			con := fmt.Sprintf("PURGE INITIATED AT %v", time.Now())

			_, err = s.ChannelMessageSendComplex(memberApprovalChannel, &discordgo.MessageSend{
				Content: con,
				Files: []*discordgo.File{
					{
						Name:        "purge.txt",
						ContentType: "text",
						Reader:      strings.NewReader(file),
					},
				},
			})
			if err != nil {
				logging.Error(s, err.Error(), i.Member.User, span, logrus.Fields{"error": err})
				return
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
		}
}
