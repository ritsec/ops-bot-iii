package commands

import (
	"fmt"

	"github.com/sirupsen/logrus"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"

	"github.com/bwmarrin/discordgo"
	"github.com/ritsec/ops-bot-iii/bot"
	"github.com/ritsec/ops-bot-iii/commands/handlers"
	"github.com/ritsec/ops-bot-iii/commands/scheduled"
	"github.com/ritsec/ops-bot-iii/commands/slash"
	"github.com/ritsec/ops-bot-iii/config"
	"github.com/ritsec/ops-bot-iii/logging"
)

var (
	// SlashCommands is a map of all slash commands
	SlashCommands map[string]func() (*discordgo.ApplicationCommand, func(s *discordgo.Session, i *discordgo.InteractionCreate)) = make(map[string]func() (*discordgo.ApplicationCommand, func(s *discordgo.Session, i *discordgo.InteractionCreate)))

	// SlashCommandHandlers is a map of all slash command handlers
	SlashCommandHandlers map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate) = make(map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate))

	// Handlers is a map of all handlers
	Handlers map[string]interface{} = make(map[string]interface{})

	// ScheduledEvents is a map of all scheduled events
	ScheduledEvents map[string]func(s *discordgo.Session, quit chan interface{}) error = make(map[string]func(s *discordgo.Session, quit chan interface{}) error)

	// ComponentHandlers is a map of all component handlers
	ComponentHandlers map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate) = make(map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate))

	// quit is a channel used to quit scheduled events
	quit chan interface{} = make(chan interface{})
)

// Init initializes all commands
func Init(ctx ddtrace.SpanContext) {
	span := tracer.StartSpan(
		"commands.Init",
		tracer.ResourceName("commands.main:Init"),
		tracer.ChildOf(ctx),
	)
	defer span.Finish()

	// Initialize the commands
	slash.Init(&ComponentHandlers, span.Context())
	scheduled.Init(&ComponentHandlers, span.Context())
	handlers.Init(&ComponentHandlers, span.Context())

	// Populate the commands
	populateSlashCommands(span.Context())
	populateHandlers(span.Context())
	populateScheduledEvents(span.Context())

	// Attach events to bot
	addSlashCommands(span.Context())
	addHandlers(span.Context())
}

// addSlashCommands adds all slash commands to the bot
func addSlashCommands(ctx ddtrace.SpanContext) {
	span := tracer.StartSpan(
		"commands.main:addSlashCommands",
		tracer.ResourceName("commands.main:addSlashCommands"),
		tracer.ChildOf(ctx),
	)
	defer span.Finish()

	// Add all slash commands to the bot
	for _, slashFunc := range SlashCommands {
		command, handler := slashFunc()
		_, err := bot.Session.ApplicationCommandCreate(config.AppID, config.GuildID, command)
		if err != nil {
			logging.Error(bot.Session, fmt.Sprintf("Failed Loading Command: %v", command.Name), nil, span, logrus.Fields{"error": err})
			logrus.Errorf("Error registered slash command: %s", command.Name)
		} else {
			logrus.Infof("Registered slash command: %s", command.Name)
			logging.Debug(bot.Session, fmt.Sprintf("Registered slash command: %s", command.Name), nil, span)
		}

		SlashCommandHandlers[command.Name] = handler
	}

	// Main handler for all events
	bot.Session.AddHandler(func(
		s *discordgo.Session,
		i *discordgo.InteractionCreate,
	) {

		switch i.Type {
		case discordgo.InteractionApplicationCommand:
			data := i.ApplicationCommandData()

			if command, ok := SlashCommandHandlers[data.Name]; ok {
				command(s, i)
			}
		case discordgo.InteractionMessageComponent:
			if command, ok := ComponentHandlers[i.MessageComponentData().CustomID]; ok {
				command(s, i)
			}
		case discordgo.InteractionModalSubmit:
			if command, ok := ComponentHandlers[i.ModalSubmitData().CustomID]; ok {
				command(s, i)
			}
		}
	})
}

// addHandlers adds all handlers to the bot
func addHandlers(ctx ddtrace.SpanContext) {
	span := tracer.StartSpan(
		"commands.main:addHandlers",
		tracer.ResourceName("commands.main:addHandlers"),
		tracer.ChildOf(ctx),
	)
	defer span.Finish()

	// Add all handlers to the bot
	for name, handler := range Handlers {
		bot.Session.AddHandler(handler)
		logrus.Infof("Registered handler: %s", name)
		logging.Debug(bot.Session, fmt.Sprintf("Registered handler: %s", name), nil, span)
	}
}

// startScheduledEvents starts all scheduled events
func StartScheduledTasks(ctx ddtrace.SpanContext) {
	span := tracer.StartSpan(
		"commands.main:StartScheduledTasks",
		tracer.ResourceName("commands.main:StartScheduledTasks"),
		tracer.ChildOf(ctx),
	)
	defer span.Finish()

	// Start all scheduled events
	for name, event := range ScheduledEvents {
		go func(eventName string, eventFunc func(s *discordgo.Session, quit chan interface{}) error) {
			for {
				err := eventFunc(bot.Session, quit)
				if err != nil {
					logging.Error(bot.Session, fmt.Sprintf("Scheduled event failed: %v", eventName), nil, span, logrus.Fields{"error": err})
				} else {
					return
				}
			}
		}(name, event)

		logrus.Infof("Starting scheduled task: %v\n", name)
		logging.Debug(bot.Session, fmt.Sprintf("Starting scheduled task: %v\n", name), nil, span)
	}
}

// stopScheduledEvents stops all scheduled events
func StopScheduledTasks(ctx ddtrace.SpanContext) {
	span := tracer.StartSpan(
		"commands.main:StopScheduledTasks",
		tracer.ResourceName("commands.main:StopScheduledTasks"),
		tracer.ChildOf(ctx),
	)
	defer span.Finish()

	// Stop all scheduled events
	if len(ScheduledEvents) > 0 {
		quit <- "sussybaka"
	}
}
