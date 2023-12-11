package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// User holds the schema definition for the User entity.
type User struct {
	ent.Schema
}

// Fields of the User.
func (User) Fields() []ent.Field {
	return []ent.Field{
		field.String("id").
			Comment("User's Discord ID").
			NotEmpty().
			Unique(),
		field.String("email").
			Comment("User's email address").
			Default(""),
		field.Int8("verification_attempts").
			Comment("Number of times the user has attempted to verify their email address").
			Default(0).
			NonNegative(),
		field.Bool("verified").
			Comment("Whether the user has been verified before").
			Default(false),
	}
}

// Edges of the User.
func (User) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("signins", Signin.Type).
			Comment("Signins made by the user"),
		edge.To("votes", Vote.Type).
			Comment("Votes made by the user"),
		edge.To("shitposts", Shitpost.Type).
			Comment("Shitposts made by the user"),
		edge.To("birthday", Birthday.Type).
			Unique().
			Comment("Birthdays of the user"),
	}
}
