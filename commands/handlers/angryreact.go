package handlers

import (
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/ritsec/ops-bot-iii/config"
	"github.com/ritsec/ops-bot-iii/helpers"
	"github.com/ritsec/ops-bot-iii/logging"
	"github.com/sirupsen/logrus"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
)

var (
	// AngryReactID is the ID of the angry react emoji
	AngryReactID = config.GetString("commands.angry_react.emoji_id")
)

// Angry Reacts to messages in #*shitpost*
func AngryReact(s *discordgo.Session, m *discordgo.MessageCreate) {

	if m.Author.ID == s.State.User.ID {
		return
	}

	channel, err := s.Channel(m.ChannelID)
	if err != nil {
		span := tracer.StartSpan(
			"commands.handlers.angryreact:AngryReact",
			tracer.ResourceName("Handlers.AngryReact"),
		)
		defer span.Finish()

		logging.Error(s, err.Error(), m.Member.User, span, logrus.Fields{"error": err})
		return
	}

	if strings.Contains(strings.ToLower(channel.Name), "shitpost") {
		span := tracer.StartSpan(
			"commands.handlers.angryreact:AngryReact",
			tracer.ResourceName("Handlers.AngryReact"),
		)
		defer span.Finish()

		err := s.MessageReactionAdd(m.ChannelID, m.ID, AngryReactID)
		if err != nil {
			logging.Error(s, err.Error(), m.Member.User, span)
		}
		logging.DebugButton(
			s,
			"Angry React Added",
			discordgo.Button{
				Label: "View Message",
				URL:   helpers.JumpURL(m.Message),
				Style: discordgo.LinkButton,
				Emoji: &discordgo.ComponentEmoji{
					Name: "👀",
				},
			},
			m.Member.User,
			span,
		)
	}
}
