package schema

import "entgo.io/ent"

// PriceSnapshot holds the schema definition for the PriceSnapshot entity.
type PriceSnapshot struct {
	ent.Schema
}

// Fields of the PriceSnapshot.
func (PriceSnapshot) Fields() []ent.Field {
	return nil
}

// Edges of the PriceSnapshot.
func (PriceSnapshot) Edges() []ent.Edge {
	return nil
}
