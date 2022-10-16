package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"github.com/google/uuid"
)

type RequestSpec map[string]interface{}

type Request struct {
	ent.Schema
}

func (Request) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).
			Default(uuid.New),
		field.String("type").NotEmpty(),
		field.String("requested_by").NotEmpty(),
		field.JSON("spec", RequestSpec{}),
	}
}

func (Request) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("Project", Project.Type).Required().Unique(),
		edge.To("Service", Service.Type).Required().Unique(),
		edge.From("Approvals", Approval.Type).Ref("Request"),
	}
}

func (Request) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("id").
			Unique(),
	}
}
