package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/field"
)

// AppSetting holds persisted application settings as key-value pairs.
type AppSetting struct {
	ent.Schema
}

func (AppSetting) Fields() []ent.Field {
	return []ent.Field{
		field.String("key").NotEmpty().Unique(),
		field.String("value").Default(""),
		field.Time("updated_at").Optional().Nillable().Default(time.Now),
	}
}

func (AppSetting) Edges() []ent.Edge {
	return nil
}
