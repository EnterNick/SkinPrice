package mixins

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect"
	"entgo.io/ent/schema/field"
	entmixin "entgo.io/ent/schema/mixin"
	"github.com/google/uuid"
)

type Identifiable struct{ entmixin.Schema }

func (Identifiable) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).
			Default(uuid.New).
			Immutable().
			Unique().
			SchemaType(map[string]string{
				dialect.Postgres: "uuid",
			}).
			StructTag(`json:"id"`),
	}
}
