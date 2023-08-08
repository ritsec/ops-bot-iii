package data

import (
	"strings"

	"gitlab.ritsec.cloud/1nv8rZim/ops-bot-iii/ent"
	"gitlab.ritsec.cloud/1nv8rZim/ops-bot-iii/ent/user"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
)

// User is the interface for interacting with the user table
type user_s struct{}

// Get gets a user by their ID, creating them if they don't exist
func (*user_s) Get(id string, ctx ddtrace.SpanContext) (*ent.User, error) {
	span := tracer.StartSpan(
		"data.user:Get",
		tracer.ResourceName("Data.User.Get"),
		tracer.ChildOf(ctx),
	)
	defer span.Finish()

	exists, err := Client.User.Query().
		Where(user.ID(id)).
		Exist(Ctx)
	if err != nil {
		return nil, err
	}

	if exists {
		return Client.User.Query().
			Where(user.ID(id)).
			Only(Ctx)
	} else {
		return Client.User.Create().
			SetID(id).
			Save(Ctx)
	}
}

// Create creates a new user
func (*user_s) Create(id string, email string, ctx ddtrace.SpanContext) (*ent.User, error) {
	span := tracer.StartSpan(
		"data.user:Create",
		tracer.ResourceName("Data.User.Create"),
		tracer.ChildOf(ctx),
	)
	defer span.Finish()

	return User.Get(id, span.Context())
}

// SetEmail sets the email for a user
func (*user_s) SetEmail(id string, email string, ctx ddtrace.SpanContext) (*ent.User, error) {
	span := tracer.StartSpan(
		"data.user:SetEmail",
		tracer.ResourceName("Data.User.SetEmail"),
		tracer.ChildOf(ctx),
	)
	defer span.Finish()

	entUser, err := User.Get(id, span.Context())
	if err != nil {
		return nil, err
	}

	email = strings.Split(email, "@")[0]

	return Client.User.UpdateOneID(entUser.ID).
		SetEmail(email).
		Save(Ctx)
}

// SetEmail sets the email for a user
func (*user_s) IncrementVerificationAttempts(id string, ctx ddtrace.SpanContext) (*ent.User, error) {
	span := tracer.StartSpan(
		"data.user:IncrementVerificationAttempts",
		tracer.ResourceName("Data.User.IncrementVerificationAttempts"),
		tracer.ChildOf(ctx),
	)
	defer span.Finish()

	entUser, err := User.Get(id, span.Context())
	if err != nil {
		return nil, err
	}

	return Client.User.UpdateOneID(entUser.ID).
		AddVerificationAttempts(1).
		Save(Ctx)
}

// SetEmail sets the email for a user
func (*user_s) MarkVerified(id string, ctx ddtrace.SpanContext) (*ent.User, error) {
	span := tracer.StartSpan(
		"data.user:MarkVerified",
		tracer.ResourceName("Data.User.MarkVerified"),
		tracer.ChildOf(ctx),
	)
	defer span.Finish()

	entUser, err := User.Get(id, span.Context())
	if err != nil {
		return nil, err
	}

	return Client.User.UpdateOneID(entUser.ID).
		SetVerified(true).
		Save(Ctx)
}

// SetEmail sets the email for a user
func (*user_s) IsVerified(id string, ctx ddtrace.SpanContext) bool {
	span := tracer.StartSpan(
		"data.user:IsVerified",
		tracer.ResourceName("Data.User.IsVerified"),
		tracer.ChildOf(ctx),
	)
	defer span.Finish()

	entUser, err := User.Get(id, span.Context())
	if err != nil {
		return false
	}

	return entUser.Verified
}

// SetEmail sets the email for a user
func (*user_s) EmailExists(id string, email string, ctx ddtrace.SpanContext) bool {
	span := tracer.StartSpan(
		"data.user:EmailExists",
		tracer.ResourceName("Data.User.EmailExists"),
		tracer.ChildOf(ctx),
	)
	defer span.Finish()

	email = strings.Split(email, "@")[0]

	exists, err := Client.User.Query().
		Where(
			user.Email(email),
			user.IDNEQ(id),
		).
		Exist(Ctx)
	if err != nil {
		return false
	}

	return exists
}

// SetEmail sets the email for a user
func (*user_s) GetVerificationAttempts(id string, ctx ddtrace.SpanContext) (int, error) {
	span := tracer.StartSpan(
		"data.user:GetVerificationAttempts",
		tracer.ResourceName("Data.User.GetVerificationAttempts"),
		tracer.ChildOf(ctx),
	)
	defer span.Finish()

	entUser, err := User.Get(id, span.Context())
	if err != nil {
		return 0, err
	}

	return int(entUser.VerificationAttempts), nil
}
