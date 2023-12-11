package slash

import (
	"fmt"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/ritsec/ops-bot-iii/commands/slash/permission"
	"github.com/ritsec/ops-bot-iii/data"
	"github.com/ritsec/ops-bot-iii/logging"
	"github.com/sirupsen/logrus"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
)

func Birthday() (*discordgo.ApplicationCommand, func(s *discordgo.Session, i *discordgo.InteractionCreate)) {
	minValue := float64(1)
	return &discordgo.ApplicationCommand{
			Name:                     "birthday",
			Description:              "Add or edit birthday",
			DefaultMemberPermissions: &permission.Member,
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionSubCommand,
					Name:        "delete",
					Description: "Delete a birthday",
				},
				{
					Type:        discordgo.ApplicationCommandOptionSubCommand,
					Name:        "add",
					Description: "Add a birthday",
					Options: []*discordgo.ApplicationCommandOption{
						{
							Type:        discordgo.ApplicationCommandOptionInteger,
							Name:        "month",
							Description: "The month of your birthday",
							Required:    true,
							Choices: []*discordgo.ApplicationCommandOptionChoice{
								{
									Name:  "January",
									Value: 1,
								},
								{
									Name:  "February",
									Value: 2,
								},
								{
									Name:  "March",
									Value: 3,
								},
								{
									Name:  "April",
									Value: 4,
								},
								{
									Name:  "May",
									Value: 5,
								},
								{
									Name:  "June",
									Value: 6,
								},
								{
									Name:  "July",
									Value: 7,
								},
								{
									Name:  "August",
									Value: 8,
								},
								{
									Name:  "September",
									Value: 9,
								},
								{
									Name:  "October",
									Value: 10,
								},
								{
									Name:  "November",
									Value: 11,
								},
								{
									Name:  "December",
									Value: 12,
								},
							},
						},
						{
							Type:        discordgo.ApplicationCommandOptionInteger,
							Name:        "day",
							Description: "The day of your birthday",
							Required:    true,
							MinValue:    &minValue,
							MaxValue:    31,
						},
					},
				},
			},
		},
		func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			span := tracer.StartSpan(
				"commands.slash.birthday:Birthday",
				tracer.ResourceName("/birthday"),
			)
			defer span.Finish()

			switch i.ApplicationCommandData().Options[0].Name {
			case "delete":
				birthdayDelete(s, i, span.Context())
			case "add":
				birthdayAdd(s, i, span.Context())
			}
		}
}

// birthdayDelete deletes a birthday from the database
func birthdayDelete(s *discordgo.Session, i *discordgo.InteractionCreate, ctx ddtrace.SpanContext) {
	span := tracer.StartSpan(
		"commands.slash.birthday:BirthdayDelete",
		tracer.ResourceName("/birthday delete"),
		tracer.ChildOf(ctx),
	)
	defer span.Finish()

	// check if birthday exists
	exists, err := data.Birthday.Exists(i.Member.User.ID, span.Context())
	if err != nil {
		logging.Error(s, "Birthday has been removed", i.Member.User, span)

		return
	}

	if exists {
		// birthday exists, remove it

		tx, err := time.LoadLocation("America/New_York")
		if err != nil {
			logging.Error(s, "encounted error when loading timezone", i.Member.User, span, logrus.Fields{"err": err.Error()})

			return
		}

		entBirthday, err := data.Birthday.Get(i.Member.User.ID, span.Context())
		if err != nil {
			logging.Error(s, "encounted error when getting birthday from user", i.Member.User, span, logrus.Fields{"err": err.Error()})

			return
		}

		// check if birthday is today
		now := time.Now().In(tx)
		if int(entBirthday.Month) == int(now.Month()) && int(entBirthday.Day) == int(now.Day()) {
			// respond to user that birthday cannot be today

			err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "You cannot remove your birthday today",
					Flags:   discordgo.MessageFlagsEphemeral,
				},
			})
			if err != nil {
				logging.Error(s, "encounted error when responding to user", i.Member.User, span, logrus.Fields{"err": err.Error()})
			}

			return
		}

		// remove birthday
		_, err = data.Birthday.Delete(i.Member.User.ID, span.Context())
		if err != nil {
			// respond to user that error occured when removing birthday
			err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "Error occured when removing birthday",
					Flags:   discordgo.MessageFlagsEphemeral,
				},
			})
			if err != nil {
				logging.Error(s, "encounted error when responding to user", i.Member.User, span, logrus.Fields{"err": err.Error()})
			}

			return
		}

		logging.Debug(s, "Birthday has been removed", i.Member.User, span)

		// respond to user that birthday has been removed
		err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Birthday has been removed",
				Flags:   discordgo.MessageFlagsEphemeral,
			},
		})
		if err != nil {
			logging.Error(s, "encounted error when responding to user", i.Member.User, span, logrus.Fields{"err": err.Error()})

			return
		}

	} else {
		// birthday does not exist

		// respond to user that birthday does not exist
		err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Birthday does not exist",
				Flags:   discordgo.MessageFlagsEphemeral,
			},
		})
		if err != nil {
			logging.Error(s, "encounted error when responding to user", i.Member.User, span, logrus.Fields{"err": err.Error()})

			return
		}

	}

}

// birthdayAdd is subcommand of birthday that adds a birthday to the database
func birthdayAdd(s *discordgo.Session, i *discordgo.InteractionCreate, ctx ddtrace.SpanContext) {
	span := tracer.StartSpan(
		"commands.slash.birthday:BirthdayAdd",
		tracer.ResourceName("/birthday add"),
		tracer.ChildOf(ctx),
	)
	defer span.Finish()

	// get user parameters
	month := int(i.ApplicationCommandData().Options[0].Options[0].IntValue())
	day := int(i.ApplicationCommandData().Options[0].Options[1].IntValue())

	tz, err := time.LoadLocation("America/New_York")
	if err != nil {
		logging.Error(s, "encounted error when loading timezone", i.Member.User, span, logrus.Fields{"err": err.Error()})

		return
	}

	// check if birthday is today
	now := time.Now().In(tz)
	if month == int(now.Month()) && day == int(now.Day()) {
		// respond to user that birthday cannot be today

		err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "You cannot set your birthday to be today",
				Flags:   discordgo.MessageFlagsEphemeral,
			},
		})
		if err != nil {
			logging.Error(s, "encounted error when responding to user", i.Member.User, span, logrus.Fields{"err": err.Error()})
		}

		return
	}

	// check is birthday already exists
	exists, err := data.Birthday.Exists(i.Member.User.ID, span.Context())
	if err != nil {
		logging.Error(s, "encounted error when checking if birthday exists", i.Member.User, span, logrus.Fields{"err": err.Error()})

		return
	}

	if exists {
		// birthday exists, update existing one

		// get current birthday
		entBirthday, err := data.Birthday.Get(i.Member.User.ID, span.Context())
		if err != nil {
			logging.Error(s, "encounted error when getting birthday from user", i.Member.User, span, logrus.Fields{"err": err.Error()})

			return
		}

		// check if birthday is today
		if int(entBirthday.Month) == int(now.Month()) && int(entBirthday.Day) == int(now.Day()) {
			// respond to user that birthday cannot be today

			err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "You cannot change your birthday today",
					Flags:   discordgo.MessageFlagsEphemeral,
				},
			})
			if err != nil {
				logging.Error(s, "encounted error when responding to user", i.Member.User, span, logrus.Fields{"err": err.Error()})
			}

			return
		}

		// update birthday
		_, err = entBirthday.Update().SetDay(day).SetMonth(month).Save(data.Ctx)
		if err != nil {
			logging.Error(s, "encounted error when updating birthday", i.Member.User, span, logrus.Fields{"err": err.Error()})

			return
		}

		// log that birthday has been updated
		logging.Debug(s, "Updated birthday for "+i.Member.User.Username, i.Member.User, span)

		// send user message that birthday was successfully updated
		err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: fmt.Sprintf("Changed birtday to %d/%d", month, day),
				Flags:   discordgo.MessageFlagsEphemeral,
			},
		})
		if err != nil {
			logging.Error(s, "encounted error when responding to user", i.Member.User, span, logrus.Fields{"err": err.Error()})

			return
		}

	} else {
		// birthday does not exist create new one

		// create new birthday for user
		_, err = data.Birthday.Create(i.Member.User.ID, day, month, span.Context())
		if err != nil {
			logging.Error(s, "encounted error when creating birthday", i.Member.User, span, logrus.Fields{"err": err.Error()})

			return
		}

		// log that birthday has been updated
		logging.Debug(s, "Created birthday for "+i.Member.User.Username, i.Member.User, span)

		// send user message that birthday was successfully updated
		err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Created birthday for " + i.Member.User.Username + " ",
				Flags:   discordgo.MessageFlagsEphemeral,
			},
		})
		if err != nil {
			logging.Error(s, "encounted error when responding to user", i.Member.User, span, logrus.Fields{"err": err.Error()})

			return
		}
	}
}
