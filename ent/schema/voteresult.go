package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
)

// VoteResult holds the schema definition for the VoteResult entity.
type VoteResult struct {
	ent.Schema
}

// Fields of the VoteResult.
func (VoteResult) Fields() []ent.Field {
	return []ent.Field{
		field.String("html").
			Comment("The vote's HTML Results"),
		field.String("plain").
			Comment("The vote's plaintext results"),
		field.String("vote_id").
			Comment("The vote's ID").
			Unique().
			NotEmpty(),
	}
}

// Edges of the VoteResult.
func (VoteResult) Edges() []ent.Edge {
	return nil
}
