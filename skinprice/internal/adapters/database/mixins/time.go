package mixins

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/dialect"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	entmixin "entgo.io/ent/schema/mixin"
)

type Timestamped struct{ entmixin.Schema }

func (Timestamped) Fields() []ent.Field {
	now := time.Now
	return []ent.Field{
		field.Time("created_at").
			Default(now).
			Immutable().
			SchemaType(map[string]string{
				dialect.Postgres: "timestamptz",
			}).
			StructTag(`json:"created_at"`),

		field.Time("updated_at").
			Default(now).
			UpdateDefault(now).
			SchemaType(map[string]string{
				dialect.Postgres: "timestamptz",
			}).
			StructTag(`json:"updated_at"`),
	}
}

type SoftDelete struct{ entmixin.Schema }

func (SoftDelete) Fields() []ent.Field {
	return []ent.Field{
		field.Time("deleted_at").
			Optional().
			Nillable().
			SchemaType(map[string]string{
				dialect.Postgres: "timestamptz",
			}).
			StructTag(`json:"deleted_at,omitempty"`),
	}
}

func (SoftDelete) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("deleted_at"),
	}
}

type Archivable struct{ entmixin.Schema }

func (Archivable) Fields() []ent.Field {
	return []ent.Field{
		field.Time("archived_at").
			Optional().
			Nillable().
			SchemaType(map[string]string{
				dialect.Postgres: "timestamptz",
			}).
			StructTag(`json:"archived_at,omitempty"`),
	}
}

func (Archivable) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("archived_at"),
	}
}
