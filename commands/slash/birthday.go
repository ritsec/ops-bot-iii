package slash

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/ritsec/ops-bot-iii/commands/slash/permission"
	"github.com/ritsec/ops-bot-iii/data"
	"github.com/ritsec/ops-bot-iii/logging"
	"github.com/sirupsen/logrus"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
)

// Birthday is a slash command that allows users to add or edit their birthday
func Birthday() (*discordgo.ApplicationCommand, func(s *discordgo.Session, i *discordgo.InteractionCreate)) {
	var minValue float64 = 1
	return &discordgo.ApplicationCommand{
			Name:                     "birthday",
			Description:              "Add or edit birthday",
			DefaultMemberPermissions: &permission.Member,
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:     discordgo.ApplicationCommandOptionInteger,
					Name:     "month",
					Required: true,
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
					Type:     discordgo.ApplicationCommandOptionInteger,
					Name:     "day",
					Required: true,
					MinValue: &minValue,
					MaxValue: 31,
				},
			},
		},
		func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			span := tracer.StartSpan(
				"commands.slash.birthday:Birthday",
				tracer.ResourceName("/birthday"),
			)
			defer span.Finish()

			// get user parameters
			month := int(i.ApplicationCommandData().Options[0].IntValue())
			day := int(i.ApplicationCommandData().Options[1].IntValue())

			// TODO: check if birthday is current date and prevent edit

			// check is birthday already exists
			exists, err := data.Birthday.Exists(i.Member.User.ID, span.Context())
			if err != nil {
				logging.Error(s, "encounted error when checking if birthday exists", i.Member.User, span, logrus.Fields{"err": err.Error()})
			}

			if exists {
				// birthday exists, update existing one

				// get current birthday
				entBirthday, err := data.Birthday.Get(i.Member.User.ID, span.Context())
				if err != nil {
					logging.Error(s, "encounted error when getting birthday from user", i.Member.User, span, logrus.Fields{"err": err.Error()})
				}

				// update birthday
				_, err = entBirthday.Update().SetDay(day).SetMonth(month).Save(data.Ctx)
				if err != nil {
					logging.Error(s, "encounted error when updating birthday", i.Member.User, span, logrus.Fields{"err": err.Error()})
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
				}

			} else {
				// birthday does not exist create new one

				// create new birthday for user
				_, err = data.Birthday.Create(i.Member.User.ID, day, month, span.Context())
				if err != nil {
					logging.Error(s, "encounted error when creating birthday", i.Member.User, span, logrus.Fields{"err": err.Error()})
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
				}
			}
		}
}

// BirthdayRemove is a slash command that allows users to remove their birthday
func BirthdayRemove() (*discordgo.ApplicationCommand, func(s *discordgo.Session, i *discordgo.InteractionCreate)) {
	return &discordgo.ApplicationCommand{
			Name:                     "birthday_remove",
			Description:              "Remove a birthday",
			DefaultMemberPermissions: &permission.Member,
		},
		func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			span := tracer.StartSpan(
				"commands.slash.birthday:BirthdayRemove",
				tracer.ResourceName("/birthday_remove"),
			)
			defer span.Finish()

			// TODO: check if birthday is current date and prevent removal

			// check if birthday exists
			exists, err := data.Birthday.Exists(i.Member.User.ID, span.Context())
			if err != nil {
				logging.Error(s, "Birthday has been removed", i.Member.User, span)
			}

			if exists {
				// birthday exists, remove it

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

					logging.Error(s, "Birthday has been removed", i.Member.User, span)

				}

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
				}
			}

		}
}
