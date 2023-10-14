package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// Shitpost holds the schema definition for the Shitposts entity.
type Shitpost struct {
	ent.Schema
}

// Fields of the Shitpost.
func (Shitpost) Fields() []ent.Field {
	return []ent.Field{
		field.String("id").
			Comment("Message ID").
			Unique().
			NotEmpty(),
		field.String("channel_id").
			Comment("Channel ID").
			NotEmpty(),
		field.Int("count").
			Comment("Shitpost Count"),
	}
}

// Edges of the Shitposts.
func (Shitpost) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("user", User.Type).
			Comment("Shitpost Author").
			Unique().
			Ref("shitposts"),
	}
}
