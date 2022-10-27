package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"entgo.io/ent/schema/mixin"
	"github.com/google/uuid"
)

type Group struct {
	ent.Schema
}

func (Group) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).
			Default(uuid.New).Unique(),
		field.String("name").NotEmpty().Unique(),
		field.Bool("is_external").
			Default(false).
			Comment("Whether the group was created via Coastline or external identity provider as part of login"),
	}
}

func (Group) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("users", User.Type).Ref("groups"),
	}
}

func (Group) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("id").
			Unique(),
	}
}

func (Group) Mixin() []ent.Mixin {
	return []ent.Mixin{
		mixin.Time{},
	}
}
