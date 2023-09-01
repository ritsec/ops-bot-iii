package data

import (
	"time"

	"gitlab.ritsec.cloud/1nv8rZim/ops-bot-iii/ent"
	"gitlab.ritsec.cloud/1nv8rZim/ops-bot-iii/ent/signin"
	"gitlab.ritsec.cloud/1nv8rZim/ops-bot-iii/ent/user"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
)

// Signin is the interface for interacting with the signin table
type signin_s struct{}

// Create creates a new signin for a user
func (*signin_s) Create(userID string, signinType signin.Type, ctx ddtrace.SpanContext) (*ent.Signin, error) {
	span := tracer.StartSpan(
		"data.signin:Create",
		tracer.ResourceName("Data.Signin.Create"),
		tracer.ChildOf(ctx),
	)
	defer span.Finish()

	entUser, err := User.Get(userID, span.Context())
	if err != nil {
		return nil, err
	}

	return Client.Signin.Create().
		SetUser(entUser).
		SetType(signinType).
		Save(Ctx)
}

// GetSignins gets all signins for a user
func (*signin_s) GetSignins(id string, ctx ddtrace.SpanContext) (int, error) {
	span := tracer.StartSpan(
		"data.signin:GetSignins",
		tracer.ResourceName("Data.Signin.GetSignins"),
		tracer.ChildOf(ctx),
	)
	defer span.Finish()

	return Client.Signin.Query().
		Where(signin.HasUserWith(user.IDEQ(id))).
		Count(Ctx)
}

// GetSigninsByType gets all signins for a user of a specific type
func (*signin_s) GetSigninsByType(id string, signinType signin.Type, ctx ddtrace.SpanContext) (int, error) {
	span := tracer.StartSpan(
		"data.signin:GetSigninsByType",
		tracer.ResourceName("Data.Signin.GetSigninsByType"),
		tracer.ChildOf(ctx),
	)
	defer span.Finish()

	return Client.Signin.Query().
		Where(
			signin.HasUserWith(user.IDEQ(id)),
			signin.TypeEQ(signinType),
		).
		Count(Ctx)
}

// RecentSignin checks if a user has signed in recently
func (*signin_s) RecentSignin(userID string, signinType signin.Type, ctx ddtrace.SpanContext) (bool, error) {
	span := tracer.StartSpan(
		"data.signin:RecentSignin",
		tracer.ResourceName("Data.Signin.RecentSignin"),
		tracer.ChildOf(ctx),
	)
	defer span.Finish()

	ok, err := Client.Signin.Query().
		Where(
			signin.HasUserWith(user.IDEQ(userID)),
			signin.TypeEQ(signinType),
			signin.TimestampGTE(time.Now().Add(-12*time.Hour)),
		).
		Exist(Ctx)
	if err != nil {
		return false, err
	}
	return ok, nil
}
