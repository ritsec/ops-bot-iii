package scheduled

import (
	"github.com/bwmarrin/discordgo"
	"github.com/robfig/cron"
	"gitlab.ritsec.cloud/1nv8rZim/ops-bot-iii/helpers"
	"gitlab.ritsec.cloud/1nv8rZim/ops-bot-iii/logging"
	"gitlab.ritsec.cloud/1nv8rZim/ops-bot-iii/structs"
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
	}
}

// checks for update every day at 2am and runs if available
func Update() *structs.ScheduledEvent {
	return structs.NewScheduledTask(
		func(s *discordgo.Session, quit chan interface{}) error {
			span := tracer.StartSpan(
				"commands.scheduled.update:Update",
				tracer.ResourceName("Scheduled.Update"),
			)
			defer span.Finish()

			c := cron.New()

			// every day at 2am
			err := c.AddFunc("0 0 2 * * *", func() { updateStatus(s, span.Context()) })
			if err != nil {
				return err
			}

			c.Start()
			<-quit
			c.Stop()

			return nil
		},
	)

}
