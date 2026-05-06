package schema

import "entgo.io/ent"

// Skin holds the schema definition for the Skin entity.
type Skin struct {
	ent.Schema
}

// Fields of the Skin.
func (Skin) Fields() []ent.Field {
	return nil
}

// Edges of the Skin.
func (Skin) Edges() []ent.Edge {
	return nil
}
