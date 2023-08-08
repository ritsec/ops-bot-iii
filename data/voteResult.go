package data

import (
	"gitlab.ritsec.cloud/1nv8rZim/ops-bot-iii/ent"
	"gitlab.ritsec.cloud/1nv8rZim/ops-bot-iii/ent/voteresult"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
)

// VoteResult is the interface for interacting with the vote_result table
type vote_result_s struct{}

// Create creates a new vote result
func (*vote_result_s) Create(voteID string, html string, plain string, ctx ddtrace.SpanContext) (*ent.VoteResult, error) {
	span := tracer.StartSpan(
		"data.voteResult:Create",
		tracer.ResourceName("data.voteResult:Create"),
		tracer.ChildOf(ctx),
	)
	defer span.Finish()

	return Client.VoteResult.Create().
		SetVoteID(voteID).
		SetHTML(html).
		SetPlain(plain).
		Save(Ctx)
}

// Get gets a vote result
func (*vote_result_s) Get(voteID string, ctx ddtrace.SpanContext) (*ent.VoteResult, error) {
	span := tracer.StartSpan(
		"data.voteResult:Get",
		tracer.ResourceName("data.voteResult:Get"),
		tracer.ChildOf(ctx),
	)
	defer span.Finish()

	return Client.VoteResult.Query().
		Where(
			voteresult.VoteIDEQ(voteID),
		).
		Only(Ctx)
}
