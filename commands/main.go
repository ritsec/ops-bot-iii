package commands

import (
	"fmt"

	"github.com/sirupsen/logrus"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"

	"github.com/bwmarrin/discordgo"
	"gitlab.ritsec.cloud/1nv8rZim/ops-bot-iii/bot"
	"gitlab.ritsec.cloud/1nv8rZim/ops-bot-iii/commands/handlers"
	"gitlab.ritsec.cloud/1nv8rZim/ops-bot-iii/commands/scheduled"
	"gitlab.ritsec.cloud/1nv8rZim/ops-bot-iii/commands/slash"
	"gitlab.ritsec.cloud/1nv8rZim/ops-bot-iii/config"
	"gitlab.ritsec.cloud/1nv8rZim/ops-bot-iii/logging"
	"gitlab.ritsec.cloud/1nv8rZim/ops-bot-iii/structs"
)

var (
	// SlashCommands is a map of all slash commands
	SlashCommands map[string]*structs.SlashCommand = make(map[string]*structs.SlashCommand)

	// Handlers is a map of all handlers
	Handlers map[string]interface{} = make(map[string]interface{})

	// ScheduledEvents is a map of all scheduled events
	ScheduledEvents map[string]*structs.ScheduledEvent = make(map[string]*structs.ScheduledEvent)

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
	for _, command := range SlashCommands {
		_, err := bot.Session.ApplicationCommandCreate(config.AppID, config.GuildID, command.Command)
		if err != nil {
			logging.Error(bot.Session, fmt.Sprintf("Failed Loading Command: %v", command.Command.Name), nil, span, logrus.Fields{"error": err})
			logrus.Errorf("Error registered slash command: %s", command.Command.Name)
		} else {
			logrus.Infof("Registered slash command: %s", command.Command.Name)
			logging.Debug(bot.Session, fmt.Sprintf("Registered slash command: %s", command.Command.Name), nil, span)
		}
	}

	// Main handler for all events
	bot.Session.AddHandler(func(
		s *discordgo.Session,
		i *discordgo.InteractionCreate,
	) {

		switch i.Type {
		case discordgo.InteractionApplicationCommand:
			data := i.ApplicationCommandData()

			if command, ok := SlashCommands[data.Name]; ok {
				command.Handler(s, i)
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
		go event.Run(bot.Session, quit)
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
