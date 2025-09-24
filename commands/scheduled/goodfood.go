package scheduled

import (
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/ritsec/ops-bot-iii/config"
	"github.com/ritsec/ops-bot-iii/helpers"
	"github.com/ritsec/ops-bot-iii/logging"
	"github.com/robfig/cron"
	"github.com/sirupsen/logrus"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
)

var (
	// GoodFoodRoleID is the ID of the Good Food role
	GoodFoodRoleID string = config.GetString("commands.good_food.role_id")

	// GoodFoodChannelID is the ID of the Good Food channel
	GoodFoodChannelID string = config.GetString("commands.good_food.channel_id")
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
			Emoji: &discordgo.ComponentEmoji{
				Name: "👀",
			},
		},
		nil,
		span,
	)
}

// GoodFood is the Good Food scheduled event
func GoodFood(s *discordgo.Session, quit chan interface{}) error {
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

	makePing := func(time string, location string) {
		must(c.AddFunc(time, func() { sendGoodFoodPing(s, location, span.Context()) }))
	}

	skipPing := func(time string) {
		must(c.AddFunc(time, func() {}))
	}

	// 11:00 AM
	makePing("0 0 11 * * MON", "RITZ")            // Monday
	makePing("0 0 11 * * TUE", "Crossroads")      // Tuesday
	makePing("0 0 11 * * WED", "Brick City Cafe") // Wednesday
	skipPing("0 0 11 * * THU")                    // Thursday
	makePing("0 0 11 * * FRI", "Brick City Cafe") // Friday
	skipPing("0 0 11 * * SAT")                    // Saturday
	skipPing("0 0 11 * * SUN")                    // Sunday

	// 4:00 PM
	skipPing("0 0 16 * * MON")               // Monday
	skipPing("0 0 16 * * TUE")               // Tuesday
	makePing("0 0 16 * * WED", "RITZ")       // Wednesday
	makePing("0 0 16 * * THU", "Crossroads") // Thursday
	skipPing("0 0 16 * * FRI")               // Friday
	skipPing("0 0 16 * * SAT")               // Saturday
	skipPing("0 0 16 * * SUN")               // Sunday

	c.Start()
	<-quit
	c.Stop()

	return nil
}
