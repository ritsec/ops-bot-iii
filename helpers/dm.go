package helpers

import (
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
