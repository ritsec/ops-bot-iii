package slash

import (
	"github.com/bwmarrin/discordgo"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
)

var (
	// ComponentHandlers is a map of all component handlers
	ComponentHandlers *map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate)
)

// Init initializes all slash commands
func Init(componentHandlers *map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate), ctx ddtrace.SpanContext) {
	span := tracer.StartSpan(
		"commands.slash.main:Init",
		tracer.ResourceName("commands.slash.main:Init"),
		tracer.ChildOf(ctx),
	)

	defer span.Finish()

	ComponentHandlers = componentHandlers
}
