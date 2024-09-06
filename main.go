package main

import (
	"os"
	"os/signal"

	"github.com/sirupsen/logrus"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"

	"github.com/ritsec/ops-bot-iii/bot"
	"github.com/ritsec/ops-bot-iii/commands"
	"github.com/ritsec/ops-bot-iii/config"
	"github.com/ritsec/ops-bot-iii/logging"
	"github.com/ritsec/ops-bot-iii/web"
)

func main() {
	// close logging file
	defer logging.Out.Close()

	// start datadog tracer
	tracer.Start(
		tracer.WithService(config.ServiceName),
		tracer.WithEnv(config.Environment),
	)
	defer tracer.Stop()

	// start main span
	span := tracer.StartSpan(
		"main",
		tracer.ResourceName("main:main"),
	)
	defer span.Finish()

	logging.InfoDD("Bot Starting", span)

	// initialize commands
	commands.Init(span.Context())

	// start bot
	err := bot.Session.Open()
	if err != nil {
		logging.CriticalDD("Bot Session Failed to open", span, logrus.Fields{"error": err})
		panic(err)
	}

	// start scheduled tasks
	commands.StartScheduledTasks(span.Context())

	logging.Debug(bot.Session, "Finished init", nil, span)

	// create stop channel
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)

	// start web server
	logrus.Println("Starting web server")
	go web.Start(span.Context())
	logging.DebugDD("Web server started", span)

	logrus.Println("Press Ctrl+C to exit")

	// wait for stop signal
	<-stop

	// stop scheduled tasks
	commands.StopScheduledTasks(span.Context())

	// stop bot
	err = bot.Session.Close()
	if err != nil {
		logging.CriticalDD("Bot Session Failed to close", span, logrus.Fields{"error": err})
		panic(err)
	}

	logrus.Println("Bot stopped")
	logging.InfoDD("Bot stopping", span)
}
