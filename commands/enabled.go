package commands

import (
	"github.com/ritsec/ops-bot-iii/commands/handlers"
	"github.com/ritsec/ops-bot-iii/commands/scheduled"
	"github.com/ritsec/ops-bot-iii/commands/slash"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
)

// populateSlashCommands populates the SlashCommands map with all of the slash commands
func populateSlashCommands(ctx ddtrace.SpanContext) {
	span := tracer.StartSpan(
		"commands.enabled:populateSlashCommands",
		tracer.ResourceName("commands.enabled:populateSlashCommands"),
		tracer.ChildOf(ctx),
	)
	defer span.Finish()

	// Populate the slash commands
	SlashCommands["ping"] = slash.Ping
	SlashCommands["purge"] = slash.Purge
	SlashCommands["kudos"] = slash.Kudos
	SlashCommands["reboot"] = slash.Reboot
	SlashCommands["log"] = slash.Log
	SlashCommands["usage"] = slash.Usage
	SlashCommands["member"] = slash.Member
	SlashCommands["signin"] = slash.Signin
	SlashCommands["vote"] = slash.Vote
	SlashCommands["feedback"] = slash.Feedback
	SlashCommands["update"] = slash.Update
	SlashCommands["query"] = slash.Query
	SlashCommands["scoreboard"] = slash.Scoreboard
	SlashCommands["birthday"] = slash.Birthday
	SlashCommands["dquery"] = slash.DQuery
}

// populateHandlers populates the Handlers map with all of the handlers
func populateHandlers(ctx ddtrace.SpanContext) {
	span := tracer.StartSpan(
		"commands.enabled:populateHandlers",
		tracer.ResourceName("commands.enabled:populateHandlers"),
		tracer.ChildOf(ctx),
	)
	defer span.Finish()

	// Populate the handlers
	Handlers["uwu"] = handlers.Uwu
	Handlers["angryReact"] = handlers.AngryReact
	Handlers["memberJoin"] = handlers.MemberJoin
	Handlers["memberLeave"] = handlers.MemberLeave
	Handlers["flag"] = handlers.Flag
	Handlers["isTheStackRunning"] = handlers.IsTheStackRunning
	Handlers["messageDelete"] = handlers.MessageDelete
	Handlers["messageEdit"] = handlers.MessageEdit
	Handlers["scoreboard"] = handlers.Scoreboard
}

// populateScheduledEvents populates the ScheduledEvents map with all of the scheduled events
func populateScheduledEvents(ctx ddtrace.SpanContext) {
	span := tracer.StartSpan(
		"commands.enabled:populateScheduledEvents",
		tracer.ResourceName("commands.enabled:populateScheduledEvents"),
		tracer.ChildOf(ctx),
	)
	defer span.Finish()

	// Populate the scheduled events
	// ScheduledEvents["goodfood"] = scheduled.GoodFood
	ScheduledEvents["heartbeat"] = scheduled.Heartbeat
	ScheduledEvents["status"] = scheduled.Status
	ScheduledEvents["update"] = scheduled.Update
	ScheduledEvents["birthday"] = scheduled.Birthday
}
