package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// Shitposts holds the schema definition for the Shitposts entity.
type Shitposts struct {
	ent.Schema
}

// Fields of the Shitposts.
func (Shitposts) Fields() []ent.Field {
	return []ent.Field{
		field.String("id").
			Comment("Message ID").
			Unique().
			NotEmpty(),
		field.Int("count").
			Comment("Shitpost Count"),
	}
}

// Edges of the Shitposts.
func (Shitposts) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("user", User.Type).
			Comment("Shitpost Author").
			Unique().
			Ref("shitposts"),
	}
}
