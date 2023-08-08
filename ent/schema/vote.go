package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// Vote holds the schema definition for the Vote entity.
type Vote struct {
	ent.Schema
}

// Fields of the Vote.
func (Vote) Fields() []ent.Field {
	return []ent.Field{
		field.String("selection").
			Comment("The user's selection").
			NotEmpty(),
		field.Int("rank").
			Comment("The selection's position").
			NonNegative(),
		field.String("vote_id").
			Comment("The vote's ID").
			NotEmpty(),
	}
}

// Edges of the Vote.
func (Vote) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("user", User.Type).
			Comment("User who voted").
			Ref("votes").
			Required().
			Unique(),
	}
}
