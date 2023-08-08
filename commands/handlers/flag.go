package handlers

import (
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/sirupsen/logrus"
	"gitlab.ritsec.cloud/1nv8rZim/ops-bot-iii/logging"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
)

// Flag is a handler that deletes messages that start with /flag
func Flag(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	if len(m.Content) >= 5 && strings.ToLower(m.Content)[1:5] == "flag" {
		span := tracer.StartSpan(
			"commands.handlers.flag:Flag",
			tracer.ResourceName("Handlers.Flag"),
		)
		defer span.Finish()

		err := s.ChannelMessageDelete(m.ChannelID, m.ID)
		if err != nil {
			logging.Error(s, err.Error(), m.Member.User, span, logrus.Fields{"error": err})
		} else {
			logging.Debug(
				s,
				"Message deleted:\n"+m.Content,
				m.Member.User,
				span,
			)
		}
	}
}
