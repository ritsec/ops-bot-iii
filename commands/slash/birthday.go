package slash

import (
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

			month := int(i.ApplicationCommandData().Options[0].IntValue())
			day := int(i.ApplicationCommandData().Options[1].IntValue())


			exists, err := data.Birthday.Exists(i.Member.User.ID, span.Context())
			if err != nil {
				logging.Error(s, "encounted error when checking if birthday exists", i.Member.User, span, logrus.Fields{"err": err.Error()})
			}

			if exists {
				entBirthday, err := data.Birthday.Get(i.Member.User.ID, span.Context())
				if err != nil {
					logging.Error(s, "encounted error when getting birthday from user", i.Member.User, span, logrus.Fields{"err": err.Error()})
				}

				_, err = entBirthday.Update().SetDay(day).SetMonth(month).Save(data.Ctx)
				if err != nil {
					logging.Error(s, "encounted error when updating birthday", i.Member.User, span, logrus.Fields{"err": err.Error()})
				}

				// Send user message that birthday was updated
				logging.Debug(s, "Updated birthday for "+i.Member.User.Username, i.Member.User, span)

				err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: "Updated birthday for " + i.Member.User.Username + " ",
						Flags:   discordgo.MessageFlagsEphemeral,
					},
				})

			} else {
				_, err = data.Birthday.Create(i.Member.User.ID, day, month, span.Context())
				if err != nil {
					logging.Error(s, "encounted error when creating birthday", i.Member.User, span, logrus.Fields{"err": err.Error()})
				}

				// send user message that birthday was created
				logging.Debug(s, "Created birthday for "+i.Member.User.Username, i.Member.User, span)

				err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: "Created birthday for " + i.Member.User.Username + " ",
						Flags:   discordgo.MessageFlagsEphemeral,
					},
				})
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
				tracer.ResourceName("/birthday remove"),
			)
			defer span.Finish()

			// code here

		}
}
