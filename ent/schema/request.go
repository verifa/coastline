package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"github.com/google/uuid"
)

type Request struct {
	ent.Schema
}

func (Request) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).
			Default(uuid.New),
		field.String("name").NotEmpty().Unique(),
	}
}

func (Request) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("Project", Project.Type).Required(),
		edge.To("Service", Service.Type).Required(),
		edge.From("Approvals", Approval.Type).Ref("Request"),
	}
}

func (Request) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("id").
			Unique(),
	}
}
