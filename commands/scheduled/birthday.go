package scheduled

import (
	"github.com/bwmarrin/discordgo"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
)

// Birthday is a scheduled event that runs at midnight to remove existing birthday roles and add new ones
func Birthday(s *discordgo.Session, quit chan interface{}) error {
	span := tracer.StartSpan(
		"commands.scheduled.birthday:Birthday",
		tracer.ResourceName("Scheduled.Birthday"),
	)
	defer span.Finish()

	// code here

	return nil
}
