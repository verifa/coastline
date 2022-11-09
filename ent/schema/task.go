package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"entgo.io/ent/schema/mixin"
	"github.com/google/uuid"
)

type Task struct {
	ent.Schema
}

func (Task) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).
			Default(uuid.New),
		field.Bytes("output"),
		field.String("error"),
	}
}

func (Task) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("Trigger", Trigger.Type).Unique().Required(),
	}
}

func (Task) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("id").
			Unique(),
	}
}

func (Task) Mixin() []ent.Mixin {
	return []ent.Mixin{
		mixin.Time{},
	}
}
