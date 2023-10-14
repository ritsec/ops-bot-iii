package data

import (
	"github.com/ritsec/ops-bot-iii/ent"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace"
)

// Shitpost is the interface for interacting with the shitpost table
type shitpost_s struct{}

// Get gets a shitpost by its ID, creating it if it doesn't exist
func (*shitpost_s) Get(id string, ctx ddtrace.SpanContext) (*ent.Shitpost, error) {}
