package biz

import (
	"context"
	"strings"
	"time"

	v1 "github.com/go-kratos/kratos-layout/api/todo/v1"

	"github.com/go-kratos/kratos/v3/errors"
	"go.einride.tech/aip/filtering"
	"go.einride.tech/aip/ordering"
)

var (
	// ErrTodoNotFound is returned when a todo does not exist.
	ErrTodoNotFound = errors.NotFound(v1.ErrorReason_TODO_NOT_FOUND.String(), "todo not found")
	// ErrTodoInvalidArgument is returned when a todo request is invalid.
	ErrTodoInvalidArgument = errors.BadRequest(v1.ErrorReason_TODO_INVALID_ARGUMENT.String(), "invalid todo argument")
)

// Todo is a Todo model.
type Todo struct {
	ID         int64
	Title      string
	Content    string
	Completed  bool
	CreateTime time.Time
	UpdateTime time.Time
}

// TodoRepo is a todo repo.
type TodoRepo interface {
	FindByID(context.Context, int64) (*Todo, error)
	ListTodos(context.Context, ...ListOption) ([]*Todo, error)
	CreateTodo(context.Context, *Todo) (*Todo, error)
	UpdateTodo(context.Context, *Todo) (*Todo, error)
	DeleteTodo(context.Context, int64) error
}

// ListOption configures todo list queries.
type ListOption func(*ListOptions)

// ListOptions are todo list query options.
type ListOptions struct {
	Filter  filtering.Filter
	OrderBy ordering.OrderBy
	Offset  int
	Limit   int
}

// ListFilter sets a standard AIP filter.
func ListFilter(filter filtering.Filter) ListOption {
	return func(o *ListOptions) {
		o.Filter = filter
	}
}

// ListOrderBy sets a standard AIP order_by value.
func ListOrderBy(orderBy ordering.OrderBy) ListOption {
	return func(o *ListOptions) {
		o.OrderBy = orderBy
	}
}

// ListOffset sets an offset.
func ListOffset(offset int) ListOption {
	return func(o *ListOptions) {
		o.Offset = offset
	}
}

// ListLimit sets a limit.
func ListLimit(limit int) ListOption {
	return func(o *ListOptions) {
		o.Limit = limit
	}
}

// TodoUsecase is a Todo usecase.
type TodoUsecase struct {
	repo TodoRepo
}

// NewTodoUsecase new a Todo usecase.
func NewTodoUsecase(repo TodoRepo) *TodoUsecase {
	return &TodoUsecase{repo: repo}
}

// GetTodo returns a todo by ID.
func (uc *TodoUsecase) GetTodo(ctx context.Context, id int64) (*Todo, error) {
	return uc.repo.FindByID(ctx, id)
}

// ListTodos lists todos.
func (uc *TodoUsecase) ListTodos(ctx context.Context, opts ...ListOption) ([]*Todo, error) {
	return uc.repo.ListTodos(ctx, opts...)
}

// CreateTodo creates a todo.
func (uc *TodoUsecase) CreateTodo(ctx context.Context, todo *Todo) (*Todo, error) {
	if err := validateTodo(todo); err != nil {
		return nil, err
	}
	return uc.repo.CreateTodo(ctx, todo)
}

// UpdateTodo updates a todo.
func (uc *TodoUsecase) UpdateTodo(ctx context.Context, todo *Todo) (*Todo, error) {
	if todo == nil || todo.ID <= 0 {
		return nil, ErrTodoInvalidArgument
	}
	if err := validateTodo(todo); err != nil {
		return nil, err
	}
	return uc.repo.UpdateTodo(ctx, todo)
}

// DeleteTodo deletes a todo.
func (uc *TodoUsecase) DeleteTodo(ctx context.Context, id int64) error {
	if id <= 0 {
		return ErrTodoInvalidArgument
	}
	return uc.repo.DeleteTodo(ctx, id)
}

func validateTodo(todo *Todo) error {
	if todo == nil || strings.TrimSpace(todo.Title) == "" {
		return ErrTodoInvalidArgument
	}
	return nil
}
