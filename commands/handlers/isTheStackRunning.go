package handlers

import (
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/ritsec/ops-bot-iii/helpers"
	"github.com/ritsec/ops-bot-iii/logging"
	"github.com/sirupsen/logrus"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
)

// IsTheStackRunning is a handler that replies to messages that contain "is the stack running?"
func IsTheStackRunning(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	if strings.Contains(strings.ToLower(m.Content), "is the stack running?") {
		span := tracer.StartSpan(
			"commands.handlers.isTheStackRunning:IsTheStackRunning",
			tracer.ResourceName("Handlers.IsTheStackRunning"),
		)
		defer span.Finish()

		message, err := s.ChannelMessageSendReply(m.ChannelID, "No :flushed:", m.Reference())
		if err != nil {
			logging.Error(s, err.Error(), m.Member.User, span, logrus.Fields{"error": err})
		} else {
			logging.DebugButton(
				s,
				"**Stack is Running**\n",
				discordgo.Button{
					Label: "View Message",
					URL:   helpers.JumpURL(message),
					Style: discordgo.LinkButton,
				},
				m.Member.User,
				span,
			)
		}
	}
}
