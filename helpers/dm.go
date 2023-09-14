package helpers

import (
	"strings"

	"github.com/bwmarrin/discordgo"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
)

// SendDirectMessage sends a direct message to a user
func SendDirectMessage(s *discordgo.Session, userID string, message string, ctx ddtrace.SpanContext) error {
	span := tracer.StartSpan(
		"helpers.dm:SendDirectMessage",
		tracer.ResourceName("Helpers.SendDirectMessage"),
		tracer.ChildOf(ctx),
	)
	defer span.Finish()

	channel, err := s.UserChannelCreate(userID)
	if err != nil {
		return err
	}
	_, err = s.ChannelMessageSend(channel.ID, message)
	return err
}

// SendDirectMessageWithFile sends a direct message to a user with a file
func SendDirectMessageWithFile(s *discordgo.Session, userID string, message string, file string, ctx ddtrace.SpanContext) error {
	span := tracer.StartSpan(
		"helpers.dm:SendDirectMessageWithFile",
		tracer.ResourceName("Helpers.SendDirectMessageWithFile"),
		tracer.ChildOf(ctx),
	)
	defer span.Finish()

	channel, err := s.UserChannelCreate(userID)
	if err != nil {
		return err
	}
	_, err = s.ChannelMessageSendComplex(channel.ID, &discordgo.MessageSend{
		Content: message,
		Files: []*discordgo.File{
			{
				Name:        "message.txt",
				ContentType: "text/plain",
				Reader:      strings.NewReader(file),
			},
		},
	})
	return err
}
