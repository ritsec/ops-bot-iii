package handlers

import (
	"fmt"
	"strings"

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

func diff(a, b string) string {
	linesA := strings.Split(a, "\n")
	linesB := strings.Split(b, "\n")

	// Generate the LCS matrix
	lcs := make([][]int, len(linesA)+1)
	for i := range lcs {
		lcs[i] = make([]int, len(linesB)+1)
	}

	for i := 1; i <= len(linesA); i++ {
		for j := 1; j <= len(linesB); j++ {
			if linesA[i-1] == linesB[j-1] {
				lcs[i][j] = lcs[i-1][j-1] + 1
			} else if lcs[i-1][j] > lcs[i][j-1] {
				lcs[i][j] = lcs[i-1][j]
			} else {
				lcs[i][j] = lcs[i][j-1]
			}
		}
	}

	// Reconstruct the diff
	i, j := len(linesA), len(linesB)
	var result []string

	for i > 0 || j > 0 {
		if i > 0 && j > 0 && linesA[i-1] == linesB[j-1] {
			i--
			j--
		} else if j > 0 && (i == 0 || lcs[i][j-1] >= lcs[i-1][j]) {
			result = append(result, fmt.Sprintf("+ %d: %s", j, linesB[j-1]))
			j--
		} else if i > 0 && (j == 0 || lcs[i][j-1] < lcs[i-1][j]) {
			result = append(result, fmt.Sprintf("- %d: %s", i, linesA[i-1]))
			i--
		}
	}

	// Reverse the result for a top-to-bottom view
	for i, j := 0, len(result)-1; i < j; i, j = i+1, j-1 {
		result[i], result[j] = result[j], result[i]
	}

	return strings.Join(result, "\n")
}

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

	message := m.BeforeDelete.Content
	message = strings.ReplaceAll(message, "```", "\\`\\`\\`")

	// https://discord.com/developers/docs/resources/channel#embed-object-embed-limits
	// 1024 characters is the max length of a field value

	// Math:
	// "```\n" + "\n```" = 8 characters => 1024 - 8 = 1016 characters or less does not need to be truncated
	// "```\n" + "...\n```" = 11 characters => 1024 - 8 = 1013 max characters if needs truncated
	if len(message) > 1016 {
		message = "```\n" + message[:1013] + "\n...```"
	} else {
		message = "```\n" + message + "\n```"
	}

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
							Value: message,
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

	messageBefore := m.BeforeUpdate.Content
	messageBefore = strings.ReplaceAll(messageBefore, "```", "\\`\\`\\`")

	messageAfter := m.Content
	messageAfter = strings.ReplaceAll(messageAfter, "```", "\\`\\`\\`")

	difference := diff(messageBefore, messageAfter)
	difference = strings.ReplaceAll(difference, "```", "\\`\\`\\`")

	// https://discord.com/developers/docs/resources/channel#embed-object-embed-limits
	// 1024 characters is the max length of a field value

	// Math:
	// "```\n" + "\n```" = 8 characters => 1024 - 8 = 1016 characters or less does not need to be truncated
	// "```\n" + "...\n```" = 11 characters => 1024 - 8 = 1013 max characters if needs truncated
	if len(messageBefore) > 1016 {
		messageBefore = "```\n" + messageBefore[:1013] + "\n...```"
	} else {
		messageBefore = "```\n" + messageBefore + "\n```"
	}

	if len(messageAfter) > 1016 {
		messageAfter = "```\n" + messageAfter[:1013] + "\n...```"
	} else {
		messageAfter = "```\n" + messageAfter + "\n```"
	}

	if len(difference) > 1016 {
		difference = "```\n" + difference[:1013] + "\n...```"
	} else {
		difference = "```\n" + difference + "\n```"
	}

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
							Value: messageAfter,
						},
						{
							Name:  "Old Message",
							Value: messageBefore,
						},
						{
							Name:  "Difference",
							Value: difference,
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
