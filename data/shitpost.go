package data

import (
	"github.com/ritsec/ops-bot-iii/ent"
	"github.com/ritsec/ops-bot-iii/ent/shitpost"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
)

// Shitpost is the interface for interacting with the shitpost table
type shitpost_s struct{}

// Get gets a shitpost by its ID
func (*shitpost_s) Get(id string, ctx ddtrace.SpanContext) (*ent.Shitpost, error) {
	span := tracer.StartSpan(
		"data.shitpost:Get",
		tracer.ResourceName("Data.Shitpost.Get"),
		tracer.ChildOf(ctx),
	)
	defer span.Finish()

	return Client.Shitpost.Query().
		Where(
			shitpost.IDEQ(id),
		).
		Only(Ctx)
}

// Create creates a new shitpost
func (*shitpost_s) Create(message_id string, user_id string, count int, ctx ddtrace.SpanContext) (*ent.Shitpost, error) {
	span := tracer.StartSpan(
		"data.shitpost:Create",
		tracer.ResourceName("Data.Shitpost.Create"),
		tracer.ChildOf(ctx),
	)
	defer span.Finish()

	ent_user, err := User.Get(user_id, span.Context())
	if err != nil {
		return nil, err
	}

	return Client.Shitpost.Create().
		SetID(message_id).
		SetUser(ent_user).
		SetCount(count).
		Save(Ctx)
}

// Delete deletes a shitpost by its ID
func (*shitpost_s) Delete(id string, ctx ddtrace.SpanContext) error {
	span := tracer.StartSpan(
		"data.shitpost:Delete",
		tracer.ResourceName("Data.Shitpost.Delete"),
		tracer.ChildOf(ctx),
	)
	defer span.Finish()

	_, err := Client.Shitpost.Delete().
		Where(
			shitpost.IDEQ(id),
		).
		Exec(Ctx)

	return err
}

// Update updates a shitpost by its ID
func (*shitpost_s) Update(message_id string, user_id string, count int, ctx ddtrace.SpanContext) (*ent.Shitpost, error) {
	span := tracer.StartSpan(
		"data.shitpost:Update",
		tracer.ResourceName("Data.Shitpost.Update"),
		tracer.ChildOf(ctx),
	)
	defer span.Finish()

	exists, err := Client.Shitpost.Query().
		Where(
			shitpost.IDEQ(message_id),
		).
		Exist(Ctx)
	if err != nil {
		return nil, err
	}

	if !exists {
		return Shitposts.Create(
			message_id,
			user_id,
			count,
			span.Context(),
		)
	}

	ent_user, err := User.Get(user_id, span.Context())
	if err != nil {
		return nil, err
	}

	return Client.Shitpost.
		UpdateOneID(message_id).
		SetUser(ent_user).
		SetCount(count).
		Save(Ctx)
}

func (*shitpost_s) GetTopXShitposts(amount int, ctx ddtrace.SpanContext) ([]*ent.Shitpost, error) {
	span := tracer.StartSpan(
		"data.shitpost:GetTopXShitposts",
		tracer.ResourceName("Data.Shitpost.GetTopXShitposts"),
		tracer.ChildOf(ctx),
	)
	defer span.Finish()

	return Client.Shitpost.Query().
		WithUser().
		Order(ent.Desc(shitpost.FieldCount)).
		Limit(amount).
		All(Ctx)
}
