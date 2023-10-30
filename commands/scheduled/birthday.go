package scheduled

import (
	"github.com/bwmarrin/discordgo"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
	"github.com/robfig/cron"
	"time"
)

// Birthday is a scheduled event that runs at midnight to remove existing birthday roles and add new ones
func Birthday(s *discordgo.Session, quit chan interface{}) error {
	span := tracer.StartSpan(
		"commands.scheduled.birthday:Birthday",
		tracer.ResourceName("Scheduled.Birthday"),
	)
	defer span.Finish()

	//Hal Williams
	est, err := time.LoadLocation("America/New_York")\
	if err != nil {
		logging.Error(s, err.Error(), nil, span)
		return err
	}

	c := cron.NewWithLocation(est)

	err = c.AddFunc("0 0 0 * * *", func() {
		//call two functions, first one to remove all roles
		//second one to add the new birthday roles
	})
	if err != nil {}
	return err
	}

	c.Start()
	<-quit
	c.Stop()


	return nil
}
