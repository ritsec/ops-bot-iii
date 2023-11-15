package scheduled

import (
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/ritsec/ops-bot-iii/config"
	"github.com/robfig/cron"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
)

var (
	birthdayRoleID string = config.GetString("commands.birthday.role_id")
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
	est, err := time.LoadLocation("America/New_York")
	if err != nil {
		logging.Error(s, err.Error(), nil, span)
		return err
	}

	c := cron.NewWithLocation(est)

	err = c.AddFunc("0 0 0 * * *", func() {
		currentTime := time.Now()
		yesterday := currentTime.add(-24 * time.Hour)
		
		//removes birthday roles
		yesterday_month := yesterday.Month()
		yesterday_day := yesterday.Day()

		entRemvBirthday, err := data.getBirthday(yesterday_month, yesterday_day, ctx)
		if err != nill {
			logging.Error(s, err.Error(), nil, span)
			return err
		}

		for _, entRemvBirthday := reange entRemvBirthday {
			removeBirthday(config.GuildID, entRemvBirthday.Edges.User.ID, birthdayRoleID)
		}

		//adds birthday roles
		current_month := currentTime.Month()
		current_day := currentTime.Day()

		entAddBirthday, err := data.getBirthday(current_month, current_day, ctx)
		if err != nill {
			logging.Error(s, err.Error(), nil, span)
			return err
		}

		for _, entAddBirthday := range entAddBirthday {
			addBirthday(config.GuildID, entAddBirthday.Edges.User.ID, "birthdayID")
		}
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
