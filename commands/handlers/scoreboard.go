package handlers

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

// Scoreboard tracks all shitpost reactions and add them for scoreboard use
func Scoreboard(s *discordgo.Session, r *discordgo.MessageReactionAdd) {
	// channelID := r.ChannelID

	fmt.Println(r.Emoji.ID)

	// channel, err := s.Channel(channelID)
	// if err != nil {
	// 	span := tracer.StartSpan(
	// 		"commands.handlers.scoreboard:Scoreboard",
	// 		tracer.ResourceName("Handlers.Scoreboard"),
	// 	)
	// 	defer span.Finish()

	// 	logging.Error(s, err.Error(), nil, span, logrus.Fields{"error": err})
	// 	return
	// }

	// if strings.Contains(strings.ToLower(channel.Name), "shitpost") {
	// 	span := tracer.StartSpan(
	// 		"commands.handlers.angryreact:AngryReact",
	// 		tracer.ResourceName("Handlers.AngryReact"),
	// 	)
	// 	defer span.Finish()

	// 	messageID := r.MessageID
	// 	message, err := s.ChannelMessage(channelID, messageID)
	// }
}
