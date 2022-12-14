package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"entgo.io/ent/schema/mixin"
	"github.com/google/uuid"
)

type Review struct {
	ent.Schema
}

func (Review) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).
			Default(uuid.New),
		field.Enum("status").Values("reject", "approve"),
		field.Enum("type").Values("user", "auto").Default("user"),
	}
}

func (Review) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("request", Request.Type).Unique().Required(),
		edge.To("created_by", User.Type).Unique().Required(),
	}
}

func (Review) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("id").
			Unique(),
	}
}

func (Review) Mixin() []ent.Mixin {
	return []ent.Mixin{
		mixin.Time{},
	}
}
