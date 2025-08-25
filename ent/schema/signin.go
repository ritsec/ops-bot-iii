package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// Signin holds the schema definition for the Signin entity.
type Signin struct {
	ent.Schema
}

// Fields of the Signin.
func (Signin) Fields() []ent.Field {
	return []ent.Field{
		field.Time("timestamp").
			Comment("Time of signin").
			Default(time.Now),
		field.Enum("type").
			Comment("Type of signin").
			Values(
				"General Meeting",
				"Red Team",
				"Red Team Recruiting",
				"Reversing",
				"RVAPT",
				"Contagion",
				"Physical",
				"Wireless",
				"IR",
				"WiCyS",
				"Zero To Hero",
				"OT Security",
				"Ops",
				"Ops IG",
				"Vulnerability Research",
				"Mentorship",
				"Other",
			),
	}
}

// Edges of the Signin.
func (Signin) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("user", User.Type).
			Comment("User who signed in").
			Ref("signins").
			Required().
			Unique(),
	}
}
