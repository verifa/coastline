package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"entgo.io/ent/schema/mixin"
	"github.com/google/uuid"
)

type Trigger struct {
	ent.Schema
}

func (Trigger) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).
			Default(uuid.New),
	}
}

func (Trigger) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("Request", Request.Type).Unique().Required(),
		edge.From("Tasks", Task.Type).Ref("Trigger"),
	}
}

func (Trigger) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("id").
			Unique(),
	}
}

func (Trigger) Mixin() []ent.Mixin {
	return []ent.Mixin{
		mixin.Time{},
	}
}
