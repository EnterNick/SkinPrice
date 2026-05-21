package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// PriceSnapshot holds the schema definition for the PriceSnapshot entity.
type PriceSnapshot struct {
	ent.Schema
}

// Fields of the PriceSnapshot.
func (PriceSnapshot) Fields() []ent.Field {
	return []ent.Field{
		field.String("market_hash_name").NotEmpty().Default(""),
		field.String("source").NotEmpty().Default(""),
		field.String("source_label").Default(""),
		field.String("page_url").Default(""),
		field.String("price_text").Default(""),
		field.Int64("price_cents").Optional().Nillable(),
		field.String("currency").Default("1"),
		field.Time("fetched_at").Default(time.Now),
		field.String("metadata").Default(""),
	}
}

func (PriceSnapshot) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("market_hash_name", "source", "fetched_at"),
		index.Fields("source", "fetched_at"),
	}
}

// Edges of the PriceSnapshot.
func (PriceSnapshot) Edges() []ent.Edge {
	return nil
}
