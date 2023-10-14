package scheduled

import (
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/ritsec/ops-bot-iii/helpers"
	"github.com/ritsec/ops-bot-iii/logging"
	"github.com/robfig/cron"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
)

// updates repo, build binary, and exists if update is available
func updateOBIII(s *discordgo.Session, ctx ddtrace.SpanContext) {
	span := tracer.StartSpan(
		"commands.scheduled.update:updateOBIII",
		tracer.ResourceName("Scheduled.Update:updateOBIII"),
		tracer.ChildOf(ctx),
	)
	defer span.Finish()

	logging.Debug(s, "Checking for update", nil, span)

	update, err := helpers.UpdateMainBranch()
	if err != nil {
		logging.Error(s, err.Error(), nil, span)
		return
	}

	if update {
		logging.Critical(s, "Update available; updating", nil, span)

		err = helpers.BuildOBIII()
		if err != nil {
			logging.Error(s, err.Error(), nil, span)
			return
		}

		err = helpers.Exit()
		if err != nil {
			logging.Error(s, err.Error(), nil, span)
			return
		}
	} else {
		logging.Debug(s, "No update available", nil, span)
	}
}

// checks for update every day at 2am and runs if available
func Update(s *discordgo.Session, quit chan interface{}) error {
	span := tracer.StartSpan(
		"commands.scheduled.update:Update",
		tracer.ResourceName("Scheduled.Update"),
	)
	defer span.Finish()

	est, err := time.LoadLocation("America/New_York")
	if err != nil {
		logging.Error(s, err.Error(), nil, span)
		return err
	}

	c := cron.NewWithLocation(est)

	// every day at 2am
	err = c.AddFunc("0 0 2 * * *", func() { updateOBIII(s, span.Context()) })
	if err != nil {
		return err
	}

	c.Start()
	<-quit
	c.Stop()

	return nil
}
