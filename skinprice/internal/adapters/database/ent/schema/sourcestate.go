package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// SourceState holds the schema definition for the SourceState entity.
type SourceState struct {
	ent.Schema
}

// Fields of the SourceState.
func (SourceState) Fields() []ent.Field {
	return []ent.Field{
		field.String("source").NotEmpty(),
		field.String("api_token_encrypted").NotEmpty(),
		field.Time("updated_at").Optional().Nillable().Default(time.Now),
	}
}

func (SourceState) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("source").Unique(),
	}
}

// Edges of the SourceState.
func (SourceState) Edges() []ent.Edge {
	return nil
}
