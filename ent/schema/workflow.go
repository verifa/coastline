package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"entgo.io/ent/schema/mixin"
	"github.com/google/uuid"
)

type Workflow struct {
	ent.Schema
}

func (Workflow) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).
			Default(uuid.New),
		field.JSON("output", Object{}),
		field.String("error"),
	}
}

func (Workflow) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("Trigger", Trigger.Type).Unique().Required(),
	}
}

func (Workflow) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("id").
			Unique(),
	}
}

func (Workflow) Mixin() []ent.Mixin {
	return []ent.Mixin{
		mixin.Time{},
	}
}
