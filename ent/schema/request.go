package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"entgo.io/ent/schema/mixin"
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
		field.String("kind").NotEmpty(),
		field.String("description").NotEmpty(),
		field.Enum("status").Values("pending_approval", "rejected", "approved").Default("pending_approval"),
		field.JSON("spec", RequestSpec{}),
	}
}

func (Request) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("project", Project.Type).Required().Unique(),
		edge.To("service", Service.Type).Required().Unique(),
		edge.To("created_by", User.Type).Required().Unique(),
		edge.From("reviews", Review.Type).Ref("request"),
		edge.From("Triggers", Trigger.Type).Ref("Request"),
	}
}

func (Request) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("id").
			Unique(),
	}
}

func (Request) Mixin() []ent.Mixin {
	return []ent.Mixin{
		mixin.Time{},
	}
}
