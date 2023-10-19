package scheduled

import (
	"github.com/bwmarrin/discordgo"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
	//Importing time to figure out when its midnight
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
	for {

		// Get the time
		currentTime := time.Now()

		// Calculate midnight
		midnight := time.Date(currentTime.Yea(), currentTime.Month(), currentTime.Day(), 0, 0, 0, 0, time.Local)
		if midnight.Before(currentTime) {
			midnight = midnight.Add(24 * time.Hour)
		}

		// Calculate the time to midnight
		timeToMidnight := midnight.Sub(currentTime)

		//	Ticker to add and remove birthday roles at midnight
		ticker := time.NewTicker(timeToMidnight)

		select {
		case <-ticker.C:
			//birthday add and remove here
		}
	}


	return nil
}
