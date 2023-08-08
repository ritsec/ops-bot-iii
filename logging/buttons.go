package logging

import (
	"fmt"
	"runtime"

	"github.com/bwmarrin/discordgo"
	"github.com/sirupsen/logrus"
	"gitlab.ritsec.cloud/1nv8rZim/ops-bot-iii/structs"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace"
)

// LogButtons logs a message with buttons
func logButtons(s *discordgo.Session, content string, buttons []discordgo.MessageComponent, user *discordgo.User, level structs.LogLevel, filename string, line int, span ddtrace.Span, fields ...logrus.Fields) *discordgo.Message {

	dd_log(level, content, filename, line, span, append(fields, logrus.Fields{"user": user})...)

	if level < LogLevel() {
		return nil
	}

	author := &discordgo.MessageEmbedAuthor{}

	if user != nil {
		author.Name = user.Username
		author.IconURL = user.AvatarURL("")
	}

	// https://github.com/bwmarrin/discordgo/wiki/FAQ#sending-embeds
	message, err := s.ChannelMessageSendComplex(
		*LevelChannelMap[level],
		&discordgo.MessageSend{
			Embeds: []*discordgo.MessageEmbed{
				{
					Author: author,
					Title:  LevelNameMap[level],
					Fields: []*discordgo.MessageEmbedField{
						{
							Name:   "File",
							Value:  filename,
							Inline: true,
						},
						{
							Name:   "Line",
							Value:  fmt.Sprint(line),
							Inline: true,
						},
						{
							Name:   "Message",
							Value:  content,
							Inline: false,
						},
					},
				},
			},
			Components: []discordgo.MessageComponent{
				discordgo.ActionsRow{
					Components: buttons,
				},
			},
		},
	)

	if err != nil {
		Error(s, err.Error(), nil, span)
	}

	return message
}

// LogButton logs a message with a button
func logButton(s *discordgo.Session, content string, button discordgo.Button, user *discordgo.User, level structs.LogLevel, filename string, line int, span ddtrace.Span, fields ...logrus.Fields) *discordgo.Message {
	return logButtons(s, content, []discordgo.MessageComponent{button}, user, level, filename, line, span, fields...)
}

// DebugLowButtons logs a message with buttons at the DebugLowLevel
func DebugLowButton(s *discordgo.Session, content string, button discordgo.Button, user *discordgo.User, span ddtrace.Span, fields ...logrus.Fields) *discordgo.Message {
	_, filename, line, _ := runtime.Caller(1)
	return logButton(s, content, button, user, DebugLowLevel, filename, line, span, fields...)
}

// DebugButtons logs a message with buttons at the DebugLevel
func DebugButton(s *discordgo.Session, content string, button discordgo.Button, user *discordgo.User, span ddtrace.Span, fields ...logrus.Fields) *discordgo.Message {
	_, filename, line, _ := runtime.Caller(1)
	return logButton(s, content, button, user, DebugLevel, filename, line, span, fields...)
}

// InfoButtons logs a message with buttons at the InfoLevel
func InfoButton(s *discordgo.Session, content string, button discordgo.Button, user *discordgo.User, span ddtrace.Span, fields ...logrus.Fields) *discordgo.Message {
	_, filename, line, _ := runtime.Caller(1)
	return logButton(s, content, button, user, InfoLevel, filename, line, span, fields...)
}

// WarningButtons logs a message with buttons at the WarningLevel
func WarningButton(s *discordgo.Session, content string, button discordgo.Button, user *discordgo.User, span ddtrace.Span, fields ...logrus.Fields) *discordgo.Message {
	_, filename, line, _ := runtime.Caller(1)
	return logButton(s, content, button, user, WarningLevel, filename, line, span, fields...)
}

// ErrorButtons logs a message with buttons at the ErrorLevel
func ErrorButton(s *discordgo.Session, content string, button discordgo.Button, user *discordgo.User, span ddtrace.Span, fields ...logrus.Fields) *discordgo.Message {
	_, filename, line, _ := runtime.Caller(1)
	return logButton(s, content, button, user, ErrorLevel, filename, line, span, fields...)
}

// CriticalButtons logs a message with buttons at the CriticalLevel
func CriticalButton(s *discordgo.Session, content string, button discordgo.Button, user *discordgo.User, span ddtrace.Span, fields ...logrus.Fields) *discordgo.Message {
	_, filename, line, _ := runtime.Caller(1)
	return logButton(s, content, button, user, CriticalLevel, filename, line, span, fields...)
}

// AuditButtons logs a message with buttons at the AuditLevel
func AuditButton(s *discordgo.Session, content string, button discordgo.Button, user *discordgo.User, span ddtrace.Span, fields ...logrus.Fields) *discordgo.Message {
	_, filename, line, _ := runtime.Caller(1)
	return logButton(s, content, button, user, AuditLevel, filename, line, span, fields...)
}

// DebugLowButtons logs a message with buttons at the DebugLowLevel
func DebugLowButtons(s *discordgo.Session, content string, buttons []discordgo.MessageComponent, user *discordgo.User, span ddtrace.Span, fields ...logrus.Fields) *discordgo.Message {
	_, filename, line, _ := runtime.Caller(1)
	return logButtons(s, content, buttons, user, DebugLowLevel, filename, line, span, fields...)
}

// DebugButtons logs a message with buttons at the DebugLevel
func DebugButtons(s *discordgo.Session, content string, buttons []discordgo.MessageComponent, user *discordgo.User, span ddtrace.Span, fields ...logrus.Fields) *discordgo.Message {
	_, filename, line, _ := runtime.Caller(1)
	return logButtons(s, content, buttons, user, DebugLevel, filename, line, span, fields...)
}

// InfoButtons logs a message with buttons at the InfoLevel
func InfoButtons(s *discordgo.Session, content string, buttons []discordgo.MessageComponent, user *discordgo.User, span ddtrace.Span, fields ...logrus.Fields) *discordgo.Message {
	_, filename, line, _ := runtime.Caller(1)
	return logButtons(s, content, buttons, user, InfoLevel, filename, line, span, fields...)
}

// WarningButtons logs a message with buttons at the WarningLevel
func WarningButtons(s *discordgo.Session, content string, buttons []discordgo.MessageComponent, user *discordgo.User, span ddtrace.Span, fields ...logrus.Fields) *discordgo.Message {
	_, filename, line, _ := runtime.Caller(1)
	return logButtons(s, content, buttons, user, WarningLevel, filename, line, span, fields...)
}

// ErrorButtons logs a message with buttons at the ErrorLevel
func ErrorButtons(s *discordgo.Session, content string, buttons []discordgo.MessageComponent, user *discordgo.User, span ddtrace.Span, fields ...logrus.Fields) *discordgo.Message {
	_, filename, line, _ := runtime.Caller(1)
	return logButtons(s, content, buttons, user, ErrorLevel, filename, line, span, fields...)
}

// CriticalButtons logs a message with buttons at the CriticalLevel
func CriticalButtons(s *discordgo.Session, content string, buttons []discordgo.MessageComponent, user *discordgo.User, span ddtrace.Span, fields ...logrus.Fields) *discordgo.Message {
	_, filename, line, _ := runtime.Caller(1)
	return logButtons(s, content, buttons, user, CriticalLevel, filename, line, span, fields...)
}

// AuditButtons logs a message with buttons at the AuditLevel
func AuditButtons(s *discordgo.Session, content string, buttons []discordgo.MessageComponent, user *discordgo.User, span ddtrace.Span, fields ...logrus.Fields) *discordgo.Message {
	_, filename, line, _ := runtime.Caller(1)
	return logButtons(s, content, buttons, user, AuditLevel, filename, line, span, fields...)
}
