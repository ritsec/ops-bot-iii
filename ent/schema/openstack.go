package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

type Openstack struct {
	ent.Schema
}

// Fields of the Openstack.
// Default is set to January 1st, 0001, 00:00:00 UTC
func (Openstack) Fields() []ent.Field {
	return []ent.Field{
		field.Time("timestamp").
			Comment("Time of last reset").
			Default(func() time.Time {
				return time.Time{}
			}),
	}
}

// Edges of the Openstack.
func (Openstack) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("user", User.Type).
			Ref("openstack").
			Unique(),
	}
}
