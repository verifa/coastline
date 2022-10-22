package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"entgo.io/ent/schema/mixin"
	"github.com/google/uuid"
)

type Service struct {
	ent.Schema
}

func (Service) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).
			Default(uuid.New),
		field.String("name").NotEmpty().Unique(),
	}
}

func (Service) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("Requests", Request.Type).Ref("Service"),
	}
}

func (Service) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("id").
			Unique(),
	}
}

func (Service) Mixin() []ent.Mixin {
	return []ent.Mixin{
		mixin.Time{},
	}
}
