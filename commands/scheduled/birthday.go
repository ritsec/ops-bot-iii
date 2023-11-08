package scheduled

import (
	"github.com/bwmarrin/discordgo"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
	"github.com/robfig/cron"
	"time"
)


func removeBirthday(GuildID, UserID, birthdayRoleID string) error {
	err := s.GuildMemberRoleRemove(GuildID, UserID, birthdayRoleID)
	if err != nil {
		logging.Error(s, err.Error(), nil, span)
		return err
	}
	return nil
}

func addBirthday(GuildID, UserID, birthdayRoleID string) error {
	err := s.GuildMemberRoleAdd(GuildID, UserID, birthdayRoleID)
	if err != nil {
		logging.Error(s, err.Error(), nil, span)
		return err
	}
	return nil
}

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
		currentTime := time.Now()

		month := currentTime.Month()
		day := currentTime.Day()

		
		//loop through all birthday users from the prior day
			//for each user run removeBirthday()
			//remove all stored birthdays
		//obtain all current birthdays
		
		entBirthday, err := data.getBirthday(month, day, ctx)
		if err != nill {
			logging.Error(s, err.Error(), nil, span)
			return err
		}

		//store all birthdays

		//loop through all new stored birthdays
		for _, entBirthday := range entBirthday {
			addBirthday(config.GuildID, entBirthday.Edges.User.ID, "birthdayID")
		}
			//for each user run addBirthday
	})
	if err != nil {
	logging.Error(s, err.Error(), nil, span)
	return err
	}

	c.Start()
	<-quit
	c.Stop()


	return nil
}
