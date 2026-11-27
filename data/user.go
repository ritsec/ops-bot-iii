package data

import (
	"fmt"
	"strings"

	"github.com/ritsec/ops-bot-iii/ent"
	"github.com/ritsec/ops-bot-iii/ent/user"
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

// IncrementVerificationAttempts increments the verification attempts after user tries to verify
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

// MarkVerified marks the user verified
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

// IsVerified checks to see if the user is verified
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

// EmailExists checks to see if the email exists for the user
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

// GetVerificationAttempts returns the amount of verification attempts the user has
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

func (*user_s) GetEmail(id string, ctx ddtrace.SpanContext) (string, error) {
	span := tracer.StartSpan(
		"data.user:GetEmail",
		tracer.ResourceName("Data.User.GetEmail"),
		tracer.ChildOf(ctx),
	)
	defer span.Finish()

	entUser, err := User.Get(id, span.Context())
	if err != nil {
		return "", err
	}

	return entUser.Email, nil
}

// FindByEmail finds the user whose email matches email. Matching order:
//  1. exact case-insensitive match after normalizing a "@g.rit.edu" domain in the input to "@rit.edu";
//  2. local-part match: the lowercased part before the first "@" with one leading "g." prefix removed, compared against the same normalization of every stored non-empty user email.
//
// Returns (nil, nil) when no user matches. Returns an error when multiple users share the matching normalized local part and no exact match exists (ambiguous).
func (*user_s) FindByEmail(email string, ctx ddtrace.SpanContext) (*ent.User, error) {
	span := tracer.StartSpan(
		"data.user:FindByEmail",
		tracer.ResourceName("Data.User.FindByEmail"),
		tracer.ChildOf(ctx),
	)
	defer span.Finish()

	normalizedEmail := normalizeRITEmail(email)

	entUsers, err := Client.User.Query().
		Where(user.EmailEqualFold(normalizedEmail)).
		All(Ctx)
	if err != nil {
		return nil, err
	}

	if len(entUsers) > 0 {
		return entUsers[0], nil
	}

	entAllUsers, err := Client.User.Query().
		Where(user.EmailNEQ("")).
		All(Ctx)
	if err != nil {
		return nil, err
	}

	var matches []*ent.User

	for _, entUser := range entAllUsers {
		if emailLocalPart(entUser.Email) == emailLocalPart(normalizedEmail) {
			matches = append(matches, entUser)
		}
	}

	switch len(matches) {
	case 0:
		return nil, nil
	case 1:
		return matches[0], nil
	default:
		return nil, fmt.Errorf("ambiguous email match for %q", email)
	}
}

// normalizeRITEmail normalizes an email address for RIT identity lookups. It
// lowercases the input and rewrites a "@g.rit.edu" domain to "@rit.edu".
func normalizeRITEmail(email string) string {
	normalized := strings.ToLower(email)
	if at := strings.LastIndex(normalized, "@"); at != -1 {
		if domain := normalized[at+1:]; domain == "g.rit.edu" {
			normalized = normalized[:at+1] + "rit.edu"
		}
	}
	return normalized
}

// emailLocalPart returns the lowercased local part of the email (the substring
// before the first "@", or the whole string if no "@" is present) with a single
// leading "g." prefix removed.
func emailLocalPart(email string) string {
	local := email
	if at := strings.Index(email, "@"); at != -1 {
		local = email[:at]
	}
	local = strings.ToLower(local)
	if strings.HasPrefix(local, "g.") {
		local = local[len("g."):]
	}
	return local
}
