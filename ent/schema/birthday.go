package schema

import "entgo.io/ent"

// Birthday holds the schema definition for the Birthday entity.
type Birthday struct {
	ent.Schema
}

// Fields of the Birthday.
func (Birthday) Fields() []ent.Field {
	return nil
}

// Edges of the Birthday.
func (Birthday) Edges() []ent.Edge {
	return nil
}
