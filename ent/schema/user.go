package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"entgo.io/ent/schema/mixin"
)

type User struct {
	ent.Schema
}

func (User) Fields() []ent.Field {
	return []ent.Field{
		field.String("sub").NotEmpty(),
		field.String("iss").NotEmpty(),
		field.String("name").NotEmpty(),
		field.String("email").Optional().Nillable(),
		field.String("picture").Optional().Nillable(),
	}
}

func (User) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("groups", Group.Type),
		edge.From("requests", Request.Type).Ref("created_by"),
		edge.From("reviews", Review.Type).Ref("created_by"),
	}
}

func (User) Indexes() []ent.Index {
	return []ent.Index{
		// User is uniquely identified by the issuer (iss) and the issuer's
		// subject identifier (sub)
		index.Fields("sub", "iss").
			Unique(),
	}
}

func (User) Mixin() []ent.Mixin {
	return []ent.Mixin{
		mixin.Time{},
	}
}
