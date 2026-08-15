package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// Todo holds the schema definition for the Todo entity.
type Todo struct {
	ent.Schema
}

// Mixin of the Todo.
func (Todo) Mixin() []ent.Mixin {
	return []ent.Mixin{
		IDMixin{},
		TimeMixin{},
	}
}

// Fields of the Todo.
func (Todo) Fields() []ent.Field {
	return []ent.Field{
		field.String("title").Default(""),
		field.String("content").Default(""),
		field.Bool("completed").Default(false),
		// status marks the row lifecycle: deletes flip it to "deleted" instead
		// of removing the row, so every read filters on "active".
		field.Enum("status").Values("active", "deleted").Default("active"),
	}
}

// Indexes of the Todo.
func (Todo) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("status"),
	}
}
