package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"github.com/google/uuid"
)

type Approval struct {
	ent.Schema
}

func (Approval) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).
			Default(uuid.New),
		field.Bool("is_automated").Default(false),
		field.String("approver").NotEmpty(),
	}
}

func (Approval) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("Request", Request.Type).Required(),
	}
}

func (Approval) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("id").
			Unique(),
	}
}
