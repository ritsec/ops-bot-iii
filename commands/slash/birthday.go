package slash

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/ritsec/ops-bot-iii/commands/slash/permission"
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

			// remove print statement (just here not to cause unused variable error)
			fmt.Println(month, day)

			// code here

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
