package data

import (
	"time"

	"github.com/ritsec/ops-bot-iii/ent"
	"github.com/ritsec/ops-bot-iii/ent/openstack"
	"github.com/ritsec/ops-bot-iii/ent/user"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
)

// Openstack is the interface for interacting with the openstack table
type openstack_s struct{}

// Exists checks if a timestamp exists for a user
func (*openstack_s) Exists(user_id string, ctx ddtrace.SpanContext) (bool, error) {
	span := tracer.StartSpan(
		"data.Openstack:Exists",
		tracer.ResourceName("Data.Openstack.Exists"),
		tracer.ChildOf(ctx),
	)
	defer span.Finish()

	return Client.Openstack.Query().
		Where(
			openstack.HasUserWith(
				user.ID(user_id),
			),
		).
		Exist(Ctx)
}

func (*openstack_s) Create(user_id string, ctx ddtrace.SpanContext) (*ent.Openstack, error) {
	span := tracer.StartSpan(
		"data.Openstack:Create",
		tracer.ResourceName("Data.Openstack.Create"),
		tracer.ChildOf(ctx),
	)
	defer span.Finish()

	entUser, err := User.Get(user_id, ctx)
	if err != nil {
		return nil, err
	}

	return Client.Openstack.Create().
		SetUser(
			entUser,
		).
		Save(Ctx)
}

// Get gets the timestamp for the last reset for Openstack account for a user
func (*openstack_s) Get(user_id string, ctx ddtrace.SpanContext) (*ent.Openstack, error) {
	span := tracer.StartSpan(
		"data.openstack:Get",
		tracer.ResourceName("Data.Openstack.Get"),
		tracer.ChildOf(ctx),
	)
	defer span.Finish()

	return Client.Openstack.Query().
		Where(
			openstack.HasUserWith(
				user.ID(user_id),
			),
		).
		WithUser().
		Only(Ctx)
}

func (*openstack_s) Update(user_id string, ctx ddtrace.SpanContext) (*ent.Openstack, error) {
	span := tracer.StartSpan(
		"data.openstack:Update",
		tracer.ResourceName("Data.Openstack.Update"),
		tracer.ChildOf(ctx),
	)
	defer span.Finish()

	user, err := Client.Openstack.Query().
		Where(
			openstack.HasUserWith(
				user.ID(user_id),
			),
		).
		Only(Ctx)
	if err != nil {
		return nil, err
	}

	tx, err := time.LoadLocation("America/New_York")
	if err != nil {
		return nil, err
	}

	updatedUser, err := Client.Openstack.UpdateOne(user).
		SetTimestamp(time.Now().In(tx)).
		Save(Ctx)
	if err != nil {
		return nil, err
	}

	return updatedUser, err
}
