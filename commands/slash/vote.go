package slash

import (
	"fmt"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/google/uuid"
	"github.com/ritsec/ops-bot-iii/commands/slash/permission"
	"github.com/ritsec/ops-bot-iii/config"
	"github.com/ritsec/ops-bot-iii/data"
	"github.com/ritsec/ops-bot-iii/helpers"
	"github.com/ritsec/ops-bot-iii/logging"
	"github.com/ritsec/ops-bot-iii/structs"
	"github.com/sirupsen/logrus"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
)

// Vote is a slash command that creates a vote
func Vote() (*discordgo.ApplicationCommand, func(s *discordgo.Session, i *discordgo.InteractionCreate)) {
	return &discordgo.ApplicationCommand{
			Name:        "vote",
			Description: "Create Vote",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "title",
					Description: "Title of the vote to be created",
					Required:    true,
				},
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "option_1",
					Description: "First option for the vote",
					Required:    true,
				},
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "option_2",
					Description: "Second option for the vote",
					Required:    true,
				},
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "option_3",
					Description: "Third option for the vote",
					Required:    false,
				},
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "option_4",
					Description: "Fourth option for the vote",
					Required:    false,
				},
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "option_5",
					Description: "Fifth option for the vote",
					Required:    false,
				},
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "option_6",
					Description: "Sixth option for the vote",
					Required:    false,
				},
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "option_7",
					Description: "Seventh option for the vote",
					Required:    false,
				},
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "option_8",
					Description: "Eighth option for the vote",
					Required:    false,
				},
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "option_9",
					Description: "Ninth option for the vote",
					Required:    false,
				},
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "option_10",
					Description: "Tenth option for the vote",
					Required:    false,
				},
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "option_11",
					Description: "Eleventh option for the vote",
					Required:    false,
				},
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "option_12",
					Description: "Twelvth option for the vote",
					Required:    false,
				},
			},
			DefaultMemberPermissions: &permission.Admin,
		},
		func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			span := tracer.StartSpan(
				"commands.slash.vote:Vote",
				tracer.ResourceName("/vote"),
			)
			defer span.Finish()

			logging.Debug(s, "Vote Command Recieved", i.Member.User, span)

			title := i.ApplicationCommandData().Options[0].StringValue()

			intOptions := len(i.ApplicationCommandData().Options) - 1
			rawOptions := func() []string {
				var options []string
				for j := 0; j < intOptions; j++ {
					options = append(options, i.ApplicationCommandData().Options[j+1].StringValue())
				}
				return options
			}()

			location, err := time.LoadLocation("America/New_York")
			if err != nil {
				logging.Error(s, err.Error(), i.Member.User, span)
			}

			voteDuration := 5 * time.Minute

			voteEndTime := time.Now().In(location).Add(voteDuration).Format("3:04PM")

			if !helpers.IsUnique(rawOptions) {
				err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Flags:   discordgo.MessageFlagsEphemeral,
						Content: "Error: All options must be unique",
					},
				})
				if err != nil {
					logging.Error(s, err.Error(), nil, span)
				}
				return
			}

			voteSlug := uuid.New().String()

			(*ComponentHandlers)[voteSlug] = func(s *discordgo.Session, i *discordgo.InteractionCreate) {
				span_voteSlug := tracer.StartSpan(
					"commands.slash.vote:Vote:voteSlug",
					tracer.ResourceName("/vote:voteSlug"),
					tracer.ChildOf(span.Context()),
				)
				defer span.Finish()

				options := make([]string, len(rawOptions))
				copy(options, rawOptions)

				signins, err := data.Signin.GetSignins(i.Member.User.ID, span.Context())
				if err != nil {
					logging.Error(s, err.Error(), i.Member.User, span_voteSlug)
					return
				}

				if signins < 3 {
					err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
						Type: discordgo.InteractionResponseChannelMessageWithSource,
						Data: &discordgo.InteractionResponseData{
							Flags:   discordgo.MessageFlagsEphemeral,
							Content: "Error: You must have at least 3 signins to vote",
						},
					})
					if err != nil {
						logging.Error(s, err.Error(), nil, span_voteSlug, logrus.Fields{"error": err})
					}
					return
				}

				i, ranking, err := voteGetVote(s, i, options, span_voteSlug.Context())
				if err != nil {
					logging.Error(s, err.Error(), nil, span_voteSlug, logrus.Fields{"error": err})
					return
				}

				for j, option := range ranking {
					_, err := data.Vote.Create(i.Member.User.ID, voteSlug, option, j, span.Context())
					if err != nil {
						logging.Error(s, err.Error(), i.Member.User, span_voteSlug, logrus.Fields{"error": err})
						return
					}
				}

				err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseUpdateMessage,
					Data: &discordgo.InteractionResponseData{
						Content: func() string {
							response := fmt.Sprintf("Voting submitted for: **%s**\n", title)
							response += "Ranking:\n"
							for _, option := range ranking {
								response += fmt.Sprintf("- **%s**\n", option)
							}
							response += "\nThanks for voting!\nIf you wish to overwrite the votes, just use the vote button to vote again!\n"
							return response
						}(),
						Components: []discordgo.MessageComponent{},
					},
				})
				if err != nil {
					logging.Error(s, err.Error(), i.Member.User, span_voteSlug, logrus.Fields{"error": err})
					return
				}
			}

			err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: func() string {
						response := fmt.Sprintf("Voting opened for: **%s**\n", title)
						response += fmt.Sprintf("Voting will close at: **%s**\n", voteEndTime)
						response += "Options:\n"
						for _, option := range rawOptions {
							response += fmt.Sprintf("- **%s**\n", option)
						}

						response += "\nVote by clicking the below button\n"

						return response
					}(),
					Components: []discordgo.MessageComponent{
						discordgo.ActionsRow{
							Components: []discordgo.MessageComponent{
								discordgo.Button{
									Label:    "Vote",
									Style:    discordgo.SuccessButton,
									CustomID: voteSlug,
								},
							},
						},
					},
				},
			})
			if err != nil {
				logging.Error(s, err.Error(), i.Member.User, span, logrus.Fields{"error": err})
				return
			}

			time.Sleep(voteDuration)

			delete(*ComponentHandlers, voteSlug)

			results, err := voteRankChoiceVoting(title, rawOptions, voteSlug, span.Context())
			if err != nil {
				logging.Error(s, err.Error(), nil, span, logrus.Fields{"error": err})
				return
			}

			entVoteResults, err := data.VoteResult.Create(voteSlug, results.HTML(), results.String(), span.Context())
			if err != nil {
				logging.Error(s, err.Error(), nil, span, logrus.Fields{"error": err})
				return
			}

			out := "```\n" + entVoteResults.Plain + "\n```" + "\nResults Graph:\n" + config.Web.Protocol + "://" + config.Web.Hostname + ":" + config.Web.Port + "/vote/" + voteSlug + "\n"

			_, err = s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
				Content:    &out,
				Components: &[]discordgo.MessageComponent{},
			})
			if err != nil {
				logging.Error(s, err.Error(), nil, span, logrus.Fields{"error": err})
				return
			}

			_, err = data.Vote.DeleteAll(voteSlug, span.Context())
			if err != nil {
				logging.Error(s, err.Error(), nil, span, logrus.Fields{"error": err})
			}
		}
}

// voteRankChoiceVoting is a helper function to calculate the results of a rank choice vote
func voteRankChoiceVoting(title string, rawOptions []string, voteID string, ctx ddtrace.SpanContext) (structs.RankChoiceVote, error) {
	span := tracer.StartSpan(
		"commands.slash.vote:voteRankChoiceVoting",
		tracer.ResourceName("/vote:rankchoice"),
		tracer.ChildOf(ctx),
	)
	defer span.Finish()

	var (
		vote    structs.RankChoiceVote
		options []string = make([]string, len(rawOptions))
	)

	vote.Title = title
	vote.Options = make([]string, len(rawOptions))

	copy(options, rawOptions)
	copy(vote.Options, rawOptions)

	for i := 1; i < len(rawOptions); i++ {
		round, err := data.Vote.GetRound(voteID, span.Context())
		if err != nil {
			return vote, err
		}

		vote.Rounds = append(vote.Rounds, structs.ConvertToRound(round, options))

		eliminationPool := helpers.EliminationPool(round, options)

		if helpers.IsEmpty(eliminationPool) {
			return vote, fmt.Errorf("empty elimination pool")
		}

		eliminated := helpers.Choose(eliminationPool)

		_, err = data.Vote.RemoveChoice(voteID, eliminated, span.Context())
		if err != nil {
			return vote, err
		}
		options = helpers.Remove(options, eliminated)
		vote.Eliminations = append(vote.Eliminations, eliminated)
	}

	round, err := data.Vote.GetRound(voteID, span.Context())
	if err != nil {
		return vote, err
	}

	vote.Rounds = append(vote.Rounds, structs.ConvertToRound(round, options))

	vote.Winner = func() string {
		for choice := range round {
			return choice
		}
		return ""
	}()

	return vote, nil
}

// voteGetVote is a helper function to get the vote from a user
func voteGetVote(s *discordgo.Session, i *discordgo.InteractionCreate, rawOptions []string, ctx ddtrace.SpanContext) (*discordgo.InteractionCreate, []string, error) {
	span := tracer.StartSpan(
		"commands.slash.vote:voteGetVote",
		tracer.ResourceName("/vote:getvote"),
		tracer.ChildOf(ctx),
	)
	defer span.Finish()

	ranking := []string{}

	options := make([]string, len(rawOptions))
	copy(options, rawOptions)

	selectionChan := make(chan string)
	defer close(selectionChan)

	interactionCreateChan := make(chan *discordgo.InteractionCreate)
	defer close(interactionCreateChan)

	messageSlug := uuid.New().String()

	(*ComponentHandlers)[messageSlug] = func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		selectionChan <- i.MessageComponentData().Values[0]
		interactionCreateChan <- i
	}

	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "Please select your highest ranked option\n",
			Flags:   discordgo.MessageFlagsEphemeral,
			Components: []discordgo.MessageComponent{
				discordgo.ActionsRow{
					Components: []discordgo.MessageComponent{
						discordgo.SelectMenu{
							CustomID: messageSlug,
							Options: func() []discordgo.SelectMenuOption {
								var innerOptions []discordgo.SelectMenuOption
								for _, option := range options {
									innerOptions = append(innerOptions, discordgo.SelectMenuOption{
										Label: option,
										Value: option,
									})
								}
								return innerOptions
							}(),
						},
					},
				},
			},
		},
	})
	if err != nil {
		return nil, nil, err
	}

	selection := <-selectionChan
	i = <-interactionCreateChan

	ranking = append(ranking, selection)
	options = helpers.Remove(options, selection)

	for len(options) > 1 {
		err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseUpdateMessage,
			Data: &discordgo.InteractionResponseData{
				Content: "Please select your highest ranked option\n",
				Flags:   discordgo.MessageFlagsEphemeral,
				Components: []discordgo.MessageComponent{
					discordgo.ActionsRow{
						Components: []discordgo.MessageComponent{
							discordgo.SelectMenu{
								CustomID: messageSlug,
								Options: func() []discordgo.SelectMenuOption {
									var innerOptions []discordgo.SelectMenuOption
									for _, option := range options {
										innerOptions = append(innerOptions, discordgo.SelectMenuOption{
											Label: option,
											Value: option,
										})
									}
									return innerOptions
								}(),
							},
						},
					},
				},
			},
		})
		if err != nil {
			return nil, nil, err
		}

		selection := <-selectionChan
		i = <-interactionCreateChan

		ranking = append(ranking, selection)
		options = helpers.Remove(options, selection)
	}

	ranking = append(ranking, options[0])

	return i, ranking, nil
}
