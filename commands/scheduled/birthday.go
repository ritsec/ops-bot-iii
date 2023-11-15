package scheduled

import (
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/ritsec/ops-bot-iii/config"
	"github.com/ritsec/ops-bot-iii/data"
	"github.com/ritsec/ops-bot-iii/helpers"
	"github.com/ritsec/ops-bot-iii/logging"
	"github.com/robfig/cron"
	"github.com/sirupsen/logrus"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
)

var (
	birthdayRoleID string = config.GetString("commands.birthday.role_id")
)

// removeBirthday removes the birthday role from a user
func removeBirthday(s *discordgo.Session, UserID string, ctx ddtrace.SpanContext) error {
	span := tracer.StartSpan(
		"commands.scheduled.birthday:removeBirthday",
		tracer.ResourceName("Scheduled.Birthday.removeBirthday"),
		tracer.ChildOf(ctx),
	)
	defer span.Finish()

	return s.GuildMemberRoleRemove(config.GuildID, UserID, birthdayRoleID)
}

// addBirthday adds the birthday role to a user
func addBirthday(s *discordgo.Session, UserID string, ctx ddtrace.SpanContext) error {
	span := tracer.StartSpan(
		"commands.scheduled.birthday:addBirthday",
		tracer.ResourceName("Scheduled.Birthday.addBirthday"),
		tracer.ChildOf(ctx),
	)
	defer span.Finish()

	return s.GuildMemberRoleAdd(config.GuildID, UserID, birthdayRoleID)
}

// Birthday is a scheduled event that runs at midnight to remove existing birthday roles and add new ones
func Birthday(s *discordgo.Session, quit chan interface{}) error {
	span := tracer.StartSpan(
		"commands.scheduled.birthday:Birthday",
		tracer.ResourceName("Scheduled.Birthday"),
	)
	defer span.Finish()

	// Set the cron job to run at EST
	est, err := time.LoadLocation("America/New_York")
	if err != nil {
		logging.Error(s, err.Error(), nil, span)
		return err
	}

	c := cron.NewWithLocation(est)

	// Run at midnight
	err = c.AddFunc("0 0 0 * * *", func() {
		internalSpan := tracer.StartSpan(
			"commands.scheduled.birthday:Birthday.Cron",
			tracer.ResourceName("Scheduled.Birthday.Cron"),
			tracer.ChildOf(span.Context()),
		)
		defer internalSpan.Finish()

		today := time.Now()
		yesterday := today.Add(-24 * time.Hour)

		// Get yesterday's birthdays
		entRemoveBirthdays, err := data.Birthday.GetBirthdays(yesterday.Day(), int(yesterday.Month()), internalSpan.Context())
		if err != nil {
			logging.Error(s, "failed to get yesterday's birthdays", nil, span, logrus.Fields{"error": err})
			return
		}

		// Remove yesterday's birthdays
		for _, entRemoveBirthday := range entRemoveBirthdays {
			err = removeBirthday(s, entRemoveBirthday.Edges.User.ID, internalSpan.Context())
			if err != nil {
				logging.Error(s, "failed to remove birthday for "+helpers.AtUser(entRemoveBirthday.Edges.User.ID), nil, span, logrus.Fields{"error": err})
				return
			}
		}

		// Get today's birthdays
		entAddBirthday, err := data.Birthday.GetBirthdays(today.Day(), int(today.Month()), internalSpan.Context())
		if err != nil {
			logging.Error(s, "failed to get today's birthdays", nil, span, logrus.Fields{"error": err})
			return
		}

		// Add today's birthdays
		for _, entAddBirthday := range entAddBirthday {
			err = addBirthday(s, entAddBirthday.Edges.User.ID, internalSpan.Context())
			if err != nil {
				logging.Error(s, "failed to add birthday for "+helpers.AtUser(entAddBirthday.Edges.User.ID), nil, span, logrus.Fields{"error": err})
				return
			}
		}
	})
	if err != nil {
		logging.Error(s, "failed to create cron job", nil, span, logrus.Fields{"error": err})
		return err
	}

	c.Start()
	<-quit
	c.Stop()

	return nil
}
