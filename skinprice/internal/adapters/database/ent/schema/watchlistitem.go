package schema

import "entgo.io/ent"

// WatchlistItem holds the schema definition for the WatchlistItem entity.
type WatchlistItem struct {
	ent.Schema
}

// Fields of the WatchlistItem.
func (WatchlistItem) Fields() []ent.Field {
	return nil
}

// Edges of the WatchlistItem.
func (WatchlistItem) Edges() []ent.Edge {
	return nil
}
