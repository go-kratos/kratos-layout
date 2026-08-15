package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"entgo.io/ent/schema/mixin"
	"github.com/google/uuid"
)

// IDMixin adds the standard application-generated UUID primary key to schemas.
type IDMixin struct {
	mixin.Schema
}

// Fields of the IDMixin.
func (IDMixin) Fields() []ent.Field {
	return []ent.Field{
		// No SchemaType: ent already maps TypeUUID per dialect (char(36) on
		// MySQL, native uuid on Postgres), so the Go side stays uuid.UUID
		// everywhere without pinning a column type.
		field.UUID("id", uuid.Nil).Default(newUUID).Immutable(),
	}
}

// newUUID returns an application-generated UUIDv7.
func newUUID() uuid.UUID {
	return uuid.Must(uuid.NewV7())
}

// TimeMixin adds standard audit timestamps to schemas.
type TimeMixin struct {
	mixin.Schema
}

// Fields of the TimeMixin.
func (TimeMixin) Fields() []ent.Field {
	return []ent.Field{
		field.Time("created_at").Immutable().Default(time.Now),
		field.Time("updated_at").Default(time.Now).UpdateDefault(time.Now),
	}
}

// Indexes of the TimeMixin.
func (TimeMixin) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("updated_at"),
	}
}
