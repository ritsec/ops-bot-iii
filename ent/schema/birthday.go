package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// Birthday holds the schema definition for the Birthday entity.
type Birthday struct {
	ent.Schema
}

// Fields of the Birthday.
func (Birthday) Fields() []ent.Field {
	return []ent.Field{
		field.Int("day").
			Comment("The day of the month of the birthday").
			Max(31).
			Positive(),
		field.Int("month").
			Comment("The month of the birthday").
			Positive().
			Max(12),
	}
}

// Edges of the Birthday.
func (Birthday) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("user", User.Type).
			Ref("birthday").
			Unique(),
	}
}
