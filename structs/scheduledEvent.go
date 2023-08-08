package structs

import (
	"github.com/sirupsen/logrus"

	"github.com/bwmarrin/discordgo"
)

// ScheduledEvent is a struct that represents a scheduled event
type ScheduledEvent struct {
	// Event is the function to run
	Event func(*discordgo.Session, chan interface{}) error
}

// NewScheduledTask creates a new scheduled task
func NewScheduledTask(function func(*discordgo.Session, chan interface{}) error) *ScheduledEvent {
	return &ScheduledEvent{
		Event: function,
	}
}

// Run runs the scheduled event
func (e *ScheduledEvent) Run(s *discordgo.Session, quit chan interface{}) {
	for {
		err := e.Event(s, quit)
		if err != nil {
			logrus.WithError(err).Error("Scheduled event failed")
		} else {
			return
		}
	}
}
