package handlers

import (
	"github.com/bwmarrin/discordgo"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
)

var (
	// ComponentHandlers is a map of component handlers
	ComponentHandlers *map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate)
)

// Init initializes the handlers
func Init(componentHandlers *map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate), ctx ddtrace.SpanContext) {
	span := tracer.StartSpan(
		"commands.handlers.main:Init",
		tracer.ResourceName("commands.handlers.main:Init"),
		tracer.ChildOf(ctx),
	)
	defer span.Finish()

	ComponentHandlers = componentHandlers
}
