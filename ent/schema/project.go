package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"entgo.io/ent/schema/mixin"
	"github.com/google/uuid"
)

type Project struct {
	ent.Schema
}

func (Project) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).
			Default(uuid.New),
		field.String("name").NotEmpty().Unique(),
	}
}

func (Project) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("requests", Request.Type).Ref("project"),
	}
}

func (Project) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("id").
			Unique(),
	}
}

func (Project) Mixin() []ent.Mixin {
	return []ent.Mixin{
		mixin.Time{},
	}
}
