package data

import (
	"github.com/ritsec/ops-bot-iii/ent"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
)

// Birthday is the interface for interacting with the birthday table
type birthday_s struct{}

// Exists checks if a birthday exists for a user
func (*birthday_s) Exists(user_id string, ctx ddtrace.SpanContext) (bool, error) {
	span := tracer.StartSpan(
		"data.birthday:Exists",
		tracer.ResourceName("Data.Birthday.Exists"),
		tracer.ChildOf(ctx),
	)
	defer span.Finish()

	// code here

	return false, nil
}

// GetBirthdays gets all birthdays for a given day and month
func (*birthday_s) GetBirthdays(day int, month int, ctx ddtrace.SpanContext) ([]ent.Birthday, error) {
	span := tracer.StartSpan(
		"data.birthday:GetBirthdays",
		tracer.ResourceName("Data.Birthday.GetBirthdays"),
		tracer.ChildOf(ctx),
	)
	defer span.Finish()

	// code here

	return nil, nil
}

// Create creates a new birthday for a user
func (*birthday_s) Create(user_id string, day int, month int, ctx ddtrace.SpanContext) (ent.Birthday, error) {
	span := tracer.StartSpan(
		"data.birthday:Create",
		tracer.ResourceName("Data.Birthday.Create"),
		tracer.ChildOf(ctx),
	)
	defer span.Finish()

	// code here

	return ent.Birthday{}, nil
}

// Get gets a birthday for a user
func (*birthday_s) Get(user_id string, ctx ddtrace.SpanContext) (ent.Birthday, error) {
	span := tracer.StartSpan(
		"data.birthday:Get",
		tracer.ResourceName("Data.Birthday.Get"),
		tracer.ChildOf(ctx),
	)
	defer span.Finish()

	// code here

	return ent.Birthday{}, nil
}

// Delete deletes a birthday for a user
func (*birthday_s) Delete(user_id string, ctx ddtrace.SpanContext) (int, error) {
	span := tracer.StartSpan(
		"data.birthday:Delete",
		tracer.ResourceName("Data.Birthday.Delete"),
		tracer.ChildOf(ctx),
	)
	defer span.Finish()

	// code here

	return 0, nil
}
