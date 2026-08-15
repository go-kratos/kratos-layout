package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"

	"github.com/go-kratos/kratos-layout/internal/biz"
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
		// status marks the row lifecycle: deletes flip it to deleted instead of
		// removing the row, so every read filters on active. Stored as an
		// integer bound to the domain type, whose values match the api enum.
		field.Int32("status").
			GoType(biz.TodoStatus(0)).
			Default(int32(biz.TodoStatusActive)),
	}
}

// Indexes of the Todo.
func (Todo) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("status"),
	}
}
