package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
)

// Skin holds the schema definition for the Skin entity.
type Skin struct {
	ent.Schema
}

// Fields of the Skin.
func (Skin) Fields() []ent.Field {
	return []ent.Field{
		field.String("market_hash_name").NotEmpty().Unique(),
		field.String("display_name").NotEmpty(),
		field.String("icon_url").Optional(),
		field.String("page_url").Optional(),
	}
}

// Edges of the Skin.
func (Skin) Edges() []ent.Edge {
	return nil
}
