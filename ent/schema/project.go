package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"github.com/google/uuid"
)

type Project struct {
	ent.Schema
}

func (Project) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "project"},
	}
}

func (Project) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).
			Default(uuid.New),
		field.String("name").NotEmpty().Unique(),
	}
}

func (Project) Edges() []ent.Edge {
	return []ent.Edge{}
}

func (Project) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("id").
			Unique(),
	}
}
