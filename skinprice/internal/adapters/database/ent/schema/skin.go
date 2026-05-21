package schema

import (
	"time"

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
		field.String("name_color").Default(""),
		field.String("icon_url").Default(""),
		field.String("page_url").Default(""),
		field.String("price_text").Default(""),
		field.String("steam_page_url").Default(""),
		field.String("steam_price_text").Default(""),
		field.Time("steam_updated_at").Optional().Nillable(),
		field.String("lisskins_page_url").Default(""),
		field.String("lisskins_price_text").Default(""),
		field.Time("lisskins_updated_at").Optional().Nillable(),
		field.String("cstm_page_url").Default(""),
		field.String("cstm_price_text").Default(""),
		field.Time("cstm_updated_at").Optional().Nillable(),
		field.String("currency").Default("1"),
		field.Time("updated_at").Optional().Nillable().Default(time.Now),
	}
}

// Edges of the Skin.
func (Skin) Edges() []ent.Edge {
	return nil
}
