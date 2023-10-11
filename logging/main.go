package logging

import (
	"fmt"
	"os"
	"runtime"

	"github.com/bwmarrin/discordgo"
	"github.com/ritsec/ops-bot-iii/config"
	"github.com/ritsec/ops-bot-iii/structs"
	"github.com/sirupsen/logrus"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace"
)

var (
	// Out is the file to write logs to
	Out *os.File
)

func init() {
	file, err := os.OpenFile(config.Logging.LogFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}

	Out = file

	logrus.SetOutput(Out)
	logrus.SetFormatter(&logrus.JSONFormatter{})
	logrus.SetLevel(logrus.TraceLevel)
}

// dd_log is a helper function to log to datadog without discord
func dd_log(level structs.LogLevel, content string, filename string, line int, span ddtrace.Span, fields ...logrus.Fields) {
	entry := logrus.WithFields(logrus.Fields{
		"span": span,
	})

	for _, field := range fields {
		entry = entry.WithFields(field)
	}

	switch level {
	case DebugLowLevel:
		entry.Debugf("[%v] %s", span, content)
	case DebugLevel:
		entry.Debugf("[%v] %s", span, content)
	case InfoLevel:
		entry.Infof("[%v] %s", span, content)
	case WarningLevel:
		entry.Warnf("[%v] %s", span, content)
	case ErrorLevel:
		entry.Errorf("[%v] %s", span, content)
	case CriticalLevel:
		entry.Fatalf("[%v] %s", span, content)
	case AuditLevel:
		entry.Infof("[%v] %s", span, content)
	}
}

// log is a helper function to log to discord and datadog
func log(s *discordgo.Session, content string, user *discordgo.User, level structs.LogLevel, filename string, line int, span ddtrace.Span, fields ...logrus.Fields) *discordgo.Message {

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
	message, _ := s.ChannelMessageSendEmbed(
		*LevelChannelMap[level],
		&discordgo.MessageEmbed{
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
	)

	return message
}

// DebugDD logs a debug message to datadog
func DebugDD(content string, span ddtrace.Span, fields ...logrus.Fields) {
	_, filename, line, _ := runtime.Caller(1)
	dd_log(DebugLevel, content, filename, line, span)
}

// InfoDD logs an info message to datadog
func InfoDD(content string, span ddtrace.Span, fields ...logrus.Fields) {
	_, filename, line, _ := runtime.Caller(1)
	dd_log(InfoLevel, content, filename, line, span)
}

// WarningDD logs a warning message to datadog
func WarningDD(content string, span ddtrace.Span, fields ...logrus.Fields) {
	_, filename, line, _ := runtime.Caller(1)
	dd_log(WarningLevel, content, filename, line, span)
}

// ErrorDD logs an error message to datadog
func ErrorDD(content string, span ddtrace.Span, fields ...logrus.Fields) {
	_, filename, line, _ := runtime.Caller(1)
	dd_log(ErrorLevel, content, filename, line, span)
}

// CriticalDD logs a critical message to datadog
func CriticalDD(content string, span ddtrace.Span, fields ...logrus.Fields) {
	_, filename, line, _ := runtime.Caller(1)
	dd_log(CriticalLevel, content, filename, line, span)
}

// AuditDD logs an audit message to datadog
func AuditDD(content string, span ddtrace.Span, fields ...logrus.Fields) {
	_, filename, line, _ := runtime.Caller(1)
	dd_log(AuditLevel, content, filename, line, span)
}

// DebugLow logs a debug message
func DebugLow(s *discordgo.Session, content string, user *discordgo.User, span ddtrace.Span, fields ...logrus.Fields) *discordgo.Message {
	_, filename, line, _ := runtime.Caller(1)
	return log(s, content, user, DebugLowLevel, filename, line, span, fields...)
}

// Debug logs a debug message
func Debug(s *discordgo.Session, content string, user *discordgo.User, span ddtrace.Span, fields ...logrus.Fields) *discordgo.Message {
	_, filename, line, _ := runtime.Caller(1)
	return log(s, content, user, DebugLevel, filename, line, span, fields...)
}

// Info logs an info message
func Info(s *discordgo.Session, content string, user *discordgo.User, span ddtrace.Span, fields ...logrus.Fields) *discordgo.Message {
	_, filename, line, _ := runtime.Caller(1)
	return log(s, content, user, InfoLevel, filename, line, span, fields...)
}

// Warning logs a warning message
func Warning(s *discordgo.Session, content string, user *discordgo.User, span ddtrace.Span, fields ...logrus.Fields) *discordgo.Message {
	_, filename, line, _ := runtime.Caller(1)
	return log(s, content, user, WarningLevel, filename, line, span, fields...)
}

// Error logs an error message
func Error(s *discordgo.Session, content string, user *discordgo.User, span ddtrace.Span, fields ...logrus.Fields) *discordgo.Message {
	_, filename, line, _ := runtime.Caller(1)
	return log(s, content, user, ErrorLevel, filename, line, span, fields...)
}

// Critical logs a critical message
func Critical(s *discordgo.Session, content string, user *discordgo.User, span ddtrace.Span, fields ...logrus.Fields) *discordgo.Message {
	_, filename, line, _ := runtime.Caller(1)
	return log(s, content, user, CriticalLevel, filename, line, span, fields...)
}

// Audit logs an audit message
func Audit(s *discordgo.Session, content string, user *discordgo.User, span ddtrace.Span, fields ...logrus.Fields) *discordgo.Message {
	_, filename, line, _ := runtime.Caller(1)
	return log(s, content, user, AuditLevel, filename, line, span, fields...)
}
