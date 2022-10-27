package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"entgo.io/ent/schema/mixin"
	"github.com/google/uuid"
)

type Session struct {
	ent.Schema
}

func (Session) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).
			Default(uuid.New),
		// Default 8 hour sessions
		field.Int64("duration").Default(int64(time.Hour * 8)),
	}
}

func (Session) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("user", User.Type).Required().Unique(),
	}
}

func (Session) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("id").
			Unique(),
	}
}

func (Session) Mixin() []ent.Mixin {
	return []ent.Mixin{
		mixin.Time{},
	}
}
