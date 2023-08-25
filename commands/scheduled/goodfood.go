package scheduled

import (
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/robfig/cron"
	"github.com/sirupsen/logrus"
	"gitlab.ritsec.cloud/1nv8rZim/ops-bot-iii/config"
	"gitlab.ritsec.cloud/1nv8rZim/ops-bot-iii/helpers"
	"gitlab.ritsec.cloud/1nv8rZim/ops-bot-iii/logging"
	"gitlab.ritsec.cloud/1nv8rZim/ops-bot-iii/structs"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
)

var (
	// GoodFoodRoleID is the ID of the Good Food role
	GoodFoodRoleID string = config.GetString("commands.goodfood.role_id")

	// GoodFoodChannelID is the ID of the Good Food channel
	GoodFoodChannelID string = config.GetString("commands.goodfood.channel_id")
)

// sendGoodFoodPing sends a ping to the Good Food channel
func sendGoodFoodPing(s *discordgo.Session, location string, ctx ddtrace.SpanContext) {
	span := tracer.StartSpan(
		"commands.scheduled.goodfood:sendGoodFoodPing",
		tracer.ResourceName("Schedued.GoodFood:sendGoodFoodPing"),
		tracer.ChildOf(ctx),
	)
	defer span.Finish()

	logging.Debug(s, "Sending Good Food ping for "+location, nil, span)

	message, err := s.ChannelMessageSend(GoodFoodChannelID, helpers.AtRole(GoodFoodRoleID)+"Good Food is Waiting @ "+location+"! https://imgur.com/ur0HqAj")
	if err != nil {
		logging.Error(s, err.Error(), message.Member.User, span, logrus.Fields{"error": err})
		return
	}
	logging.DebugButton(
		s,
		"Good Food is Here!",
		discordgo.Button{
			Label: "View Message",
			URL:   helpers.JumpURL(message),
			Style: discordgo.LinkButton,
		},
		nil,
		span,
	)
}

// GoodFood is the Good Food scheduled event
func GoodFood() *structs.ScheduledEvent {
	return structs.NewScheduledTask(
		func(s *discordgo.Session, quit chan interface{}) error {
			span := tracer.StartSpan(
				"commands.scheduled.goodfood:GoodFood",
				tracer.ResourceName("Scheduled.GoodFood"),
			)
			defer span.Finish()

			est, err := time.LoadLocation("America/New_York")
			if err != nil {
				logging.Error(s, err.Error(), nil, span)
				return err
			}

			c := cron.NewWithLocation(est)

			must := func(err error) {
				if err != nil {
					logging.Error(s, err.Error(), nil, span)
				}
			}

			// 11:00 AM
			must(c.AddFunc("0 0 11 * * MON", func() {}))                                                    // Monday
			must(c.AddFunc("0 0 11 * * TUE", func() { sendGoodFoodPing(s, "Crossroads", span.Context()) })) // Tuesday
			must(c.AddFunc("0 0 11 * * WED", func() { sendGoodFoodPing(s, "Brick City", span.Context()) })) // Wednesday
			must(c.AddFunc("0 0 11 * * THU", func() {}))                                                    // Thursday
			must(c.AddFunc("0 0 11 * * FRI", func() {}))                                                    // Friday
			must(c.AddFunc("0 0 11 * * SAT", func() {}))                                                    // Saturday
			must(c.AddFunc("0 0 11 * * SUN", func() {}))                                                    // Sunday

			// 4:00 PM
			must(c.AddFunc("0 0 16 * * MON", func() { sendGoodFoodPing(s, "RITZ", span.Context()) }))       // Monday
			must(c.AddFunc("0 0 16 * * TUE", func() {}))                                                    // Tuesday
			must(c.AddFunc("0 0 16 * * WED", func() { sendGoodFoodPing(s, "RITZ", span.Context()) }))       // Wednesday
			must(c.AddFunc("0 0 16 * * THU", func() { sendGoodFoodPing(s, "Crossroads", span.Context()) })) // Thursday
			must(c.AddFunc("0 0 16 * * FRI", func() {}))                                                    // Friday
			must(c.AddFunc("0 0 16 * * SAT", func() {}))                                                    // Saturday
			must(c.AddFunc("0 0 16 * * SUN", func() {}))                                                    // Sunday

			c.Start()
			<-quit
			c.Stop()

			return nil
		},
	)

}
