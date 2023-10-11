package data

import (
	"github.com/ritsec/ops-bot-iii/ent"
	"github.com/ritsec/ops-bot-iii/ent/user"
	"github.com/ritsec/ops-bot-iii/ent/vote"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
)

// Vote is the interface for interacting with the vote table
type vote_s struct{}

// Create creates a new vote
func (*vote_s) Create(userID string, voteID string, choice string, rank int, ctx ddtrace.SpanContext) (*ent.Vote, error) {
	span := tracer.StartSpan(
		"data.vote:Create",
		tracer.ResourceName("Data.Vote.Create"),
		tracer.ChildOf(ctx),
	)
	defer span.Finish()

	entUser, err := User.Get(userID, span.Context())
	if err != nil {
		return nil, err
	}

	exists, err := Client.Vote.Query().
		Where(
			vote.HasUserWith(user.ID(entUser.ID)),
			vote.VoteIDEQ(voteID),
			vote.Rank(rank),
		).
		Exist(Ctx)
	if err != nil {
		return nil, err
	}

	if exists {
		entVote, err := Client.Vote.Query().
			Where(
				vote.HasUserWith(user.ID(userID)),
				vote.VoteIDEQ(voteID),
				vote.Rank(rank),
			).
			Only(Ctx)
		if err != nil {
			return nil, err
		}

		entVote, err = Client.Vote.UpdateOne(entVote).
			SetSelection(choice).
			Save(Ctx)
		if err != nil {
			return nil, err
		}

		return entVote, nil
	}

	return Client.Vote.Create().
		SetUser(entUser).
		SetVoteID(voteID).
		SetSelection(choice).
		SetRank(rank).
		Save(Ctx)
}

// GetRound gets the current round of voting
func (*vote_s) GetRound(voteID string, ctx ddtrace.SpanContext) (map[string]int, error) {
	span := tracer.StartSpan(
		"data.vote:GetRound",
		tracer.ResourceName("Data.Vote.GetRound"),
		tracer.ChildOf(ctx),
	)
	defer span.Finish()

	entUsers, err := Client.User.Query().
		Where(
			user.HasVotesWith(vote.VoteIDEQ(voteID)),
		).
		All(Ctx)
	if err != nil {
		return nil, err
	}

	round := make(map[string]int)

	for _, entUser := range entUsers {
		entVote, err := Client.Vote.Query().
			Where(
				vote.HasUserWith(user.ID(entUser.ID)),
				vote.VoteIDEQ(voteID),
			).
			Order(ent.Asc(vote.FieldRank)).
			First(Ctx)
		if err != nil {
			return nil, err
		}

		round[entVote.Selection]++
	}

	return round, nil
}

// RemoveChoice removes a choice from the vote
func (*vote_s) RemoveChoice(voteID string, choice string, ctx ddtrace.SpanContext) (int, error) {
	span := tracer.StartSpan(
		"data.vote:RemoveChoice",
		tracer.ResourceName("Data.Vote.RemoveChoice"),
		tracer.ChildOf(ctx),
	)
	defer span.Finish()

	return Client.Vote.Delete().
		Where(
			vote.VoteIDEQ(voteID),
			vote.SelectionEQ(choice),
		).
		Exec(Ctx)
}

// DeleteAll deletes all votes for a given voteID
func (*vote_s) DeleteAll(voteID string, ctx ddtrace.SpanContext) (int, error) {
	span := tracer.StartSpan(
		"data.vote:DeleteAll",
		tracer.ResourceName("Data.Vote.DeleteAll"),
		tracer.ChildOf(ctx),
	)
	defer span.Finish()

	return Client.Vote.Delete().
		Where(vote.VoteIDEQ(voteID)).
		Exec(Ctx)
}
