package slash

import (
	"fmt"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/ritsec/ops-bot-iii/commands/slash/permission"
	"github.com/ritsec/ops-bot-iii/config"
	"github.com/ritsec/ops-bot-iii/logging"
	"github.com/sirupsen/logrus"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
)

var (
	// Channel ID of the purge logs channel
	purgeLogsChannel string = config.GetString("commands.purge.channel_id")
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
				timeloc     *time.Location
			)

			timeloc, err := time.LoadLocation("EST")

			if err != nil {
				logging.Error(s, err.Error(), i.Member.User, span, logrus.Fields{"error": err})
			}

			if len(i.ApplicationCommandData().Options) != 0 {
				messages = i.ApplicationCommandData().Options[0].IntValue()
			} else {
				messages = 100
			}

			raw_messages, err := s.ChannelMessages(i.ChannelID, int(messages), "", "", "")
			if err != nil {
				logging.Error(s, err.Error(), i.Member.User, span, logrus.Fields{"error": err})
			}

			file := fmt.Sprintf("Record of the purge on %v\n", time.Now().In(timeloc).Format(" 15:04:05"))
			file += "Purged " + fmt.Sprint(len(raw_messages)) + " messages!\n"
			file += "------------------------------------------------"

			// For the file
			// reverses the list of messages to make the file from oldest to newest
			reversedMessages := make([]*discordgo.Message, len(raw_messages))
			for j, message := range raw_messages {
				reversedMessages[len(raw_messages)-1-j] = message
			}

			for _, message := range reversedMessages {
				// Check to see if message is edited
				if message.EditedTimestamp == nil {
					file += fmt.Sprintf("\n\n%v SENT AT %v", message.Author, message.Timestamp.In(timeloc).Format("2006-01-02 15:04:05"))
					file += fmt.Sprintf("\n%v", message.Content)
				} else {
					file += fmt.Sprintf("\n\n%v SENT AT %v (EDITED AT %v)", message.Author, message.Timestamp.In(timeloc).Format("2006-01-02 15:04:05"), message.EditedTimestamp.In(timeloc).Format("2006-01-02 15:04:05"))
					file += fmt.Sprintf("\n%v", message.Content)
				}

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

			// Putting this here so that it does not send the file if the ChannelMessagesBulkDelete fails for some reason

			con := fmt.Sprintf("PURGE INITIATED AT %v", time.Now().Local().Format("2006-01-02 15:04:05-07:00"))
			con += "Purged " + fmt.Sprint(len(raw_messages)) + " messages!"

			_, err = s.ChannelMessageSendComplex(purgeLogsChannel, &discordgo.MessageSend{
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
		}
}
