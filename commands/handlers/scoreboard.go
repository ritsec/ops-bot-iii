package handlers

import (
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/ritsec/ops-bot-iii/config"
	"github.com/ritsec/ops-bot-iii/data"
	"github.com/ritsec/ops-bot-iii/logging"
	"github.com/sirupsen/logrus"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
)

var (
	// AngryReactID is the ID of the angry react emoji
	AngryReactEmojiID = strings.Split(config.GetString("commands.angry_react.emoji_id"), ":")[1]
)

// Scoreboard tracks all shitpost reactions and add them for scoreboard use
func Scoreboard(s *discordgo.Session, r *discordgo.MessageReactionAdd) {
	if r.Emoji.ID != AngryReactEmojiID {
		return
	}

	channelID := r.ChannelID
	channel, err := s.Channel(channelID)
	if err != nil {
		span := tracer.StartSpan(
			"commands.handlers.scoreboard:Scoreboard",
			tracer.ResourceName("Handlers.Scoreboard"),
		)
		defer span.Finish()

		logging.Error(s, err.Error(), nil, span, logrus.Fields{"error": err})
		return
	}

	if strings.Contains(strings.ToLower(channel.Name), "shitpost") {
		span := tracer.StartSpan(
			"commands.handlers.angryreact:AngryReact",
			tracer.ResourceName("Handlers.AngryReact"),
		)
		defer span.Finish()

		message, err := s.ChannelMessage(channelID, r.MessageID)
		if err != nil {
			logging.Error(s, err.Error(), nil, span)
			return
		}

		var count int
		for _, reaction := range message.Reactions {
			if reaction.Emoji.ID == AngryReactEmojiID {
				count = reaction.Count
			}
		}

		_, err = data.Shitposts.Update(message.ID, message.ChannelID, r.Member.User.ID, count, span.Context())
		if err != nil {
			logging.Error(s, err.Error(), nil, span)
			return
		}
	}
}
