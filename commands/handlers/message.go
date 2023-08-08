package handlers

import (
	"github.com/bwmarrin/discordgo"
	"github.com/sirupsen/logrus"
	"gitlab.ritsec.cloud/1nv8rZim/ops-bot-iii/config"
	"gitlab.ritsec.cloud/1nv8rZim/ops-bot-iii/logging"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
)

var (
	// messageModificationChannelID is the channel ID of the channel to send message modification events to
	messageModificationChannelID string = config.GetString("commands.message.channel_id")
)

// MessageDelete is a handler that sends a message to the messageModificationChannelID channel when a message is deleted
func MessageDelete(s *discordgo.Session, m *discordgo.MessageDelete) {
	if m.BeforeDelete == nil {
		return
	}

	if m.BeforeDelete.Author.Bot {
		return
	}

	author := &discordgo.MessageEmbedAuthor{}

	if m.BeforeDelete.Author != nil {
		author.Name = m.BeforeDelete.Author.Username
		author.IconURL = m.BeforeDelete.Author.AvatarURL("")
	}

	span := tracer.StartSpan(
		"commands.handlers.message:MessageDelete",
		tracer.ResourceName("Handlers.MessageDelete"),
	)
	defer span.Finish()

	_, err := s.ChannelMessageSendComplex(
		messageModificationChannelID,
		&discordgo.MessageSend{
			Embeds: []*discordgo.MessageEmbed{
				{
					Author: author,
					Title:  "Message Deleted",
					Fields: []*discordgo.MessageEmbedField{
						{
							Name:  "Channel",
							Value: "<#" + m.BeforeDelete.ChannelID + ">",
						},
						{
							Name:  "Message",
							Value: m.BeforeDelete.Content,
						},
					},
				},
			},
		},
	)
	if err != nil {
		logging.Error(s, err.Error(), nil, span, logrus.Fields{"error": err})
	}
}

// MessageEdit is a handler that sends a message to the messageModificationChannelID channel when a message is editted
func MessageEdit(s *discordgo.Session, m *discordgo.MessageUpdate) {
	if m.BeforeUpdate == nil {
		return
	}

	if m.BeforeUpdate.Author.Bot {
		return
	}

	author := &discordgo.MessageEmbedAuthor{}

	if m.BeforeUpdate.Author != nil {
		author.Name = m.BeforeUpdate.Author.Username
		author.IconURL = m.BeforeUpdate.Author.AvatarURL("")
	}

	span := tracer.StartSpan(
		"commands.handlers.message:MessageEdit",
		tracer.ResourceName("Handlers.MessageEdit"),
	)
	defer span.Finish()

	_, err := s.ChannelMessageSendComplex(
		messageModificationChannelID,
		&discordgo.MessageSend{
			Embeds: []*discordgo.MessageEmbed{
				{
					Author: author,
					Title:  "Message Editted",
					Fields: []*discordgo.MessageEmbedField{
						{
							Name:  "Channel",
							Value: "<#" + m.BeforeUpdate.ChannelID + ">",
						},
						{
							Name:  "Editted Message",
							Value: m.Content,
						},
						{
							Name:  "Old Message",
							Value: m.BeforeUpdate.Content,
						},
					},
				},
			},
		},
	)
	if err != nil {
		logging.Error(s, err.Error(), nil, span, logrus.Fields{"error": err})
	}
}
