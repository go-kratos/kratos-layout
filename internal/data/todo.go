package data

import (
	"cmp"
	"context"
	"slices"
	"sync"
	"time"

	"github.com/go-kratos/kratos-layout/internal/biz"
)

type todoRepo struct {
	data *Data

	mu     sync.RWMutex
	nextID int64
	todos  map[int64]*biz.Todo
}

// NewTodoRepo creates a new TodoRepo instance.
func NewTodoRepo(data *Data) biz.TodoRepo {
	return &todoRepo{
		data:   data,
		nextID: 1,
		todos:  make(map[int64]*biz.Todo),
	}
}

func (r *todoRepo) FindByID(_ context.Context, id int64) (*biz.Todo, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	todo, ok := r.todos[id]
	if !ok {
		return nil, biz.ErrTodoNotFound
	}
	return cloneTodo(todo), nil
}

func (r *todoRepo) ListTodos(_ context.Context, opts ...biz.ListOption) ([]*biz.Todo, error) {
	options := biz.ListOptions{Limit: 20}
	for _, opt := range opts {
		opt(&options)
	}
	if options.Offset < 0 || options.Limit <= 0 {
		return nil, biz.ErrTodoInvalidArgument
	}

	r.mu.RLock()
	defer r.mu.RUnlock()

	todos := make([]*biz.Todo, 0, len(r.todos))
	for _, todo := range r.todos {
		todos = append(todos, cloneTodo(todo))
	}
	slices.SortFunc(todos, func(a, b *biz.Todo) int {
		return cmp.Compare(a.ID, b.ID)
	})

	if options.Offset >= len(todos) {
		return []*biz.Todo{}, nil
	}
	end := options.Offset + options.Limit
	if end > len(todos) {
		end = len(todos)
	}
	return todos[options.Offset:end], nil
}

func (r *todoRepo) CreateTodo(_ context.Context, todo *biz.Todo) (*biz.Todo, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	now := time.Now()
	todo = cloneTodo(todo)
	todo.ID = r.nextID
	todo.CreateTime = now
	todo.UpdateTime = now
	r.todos[todo.ID] = cloneTodo(todo)
	r.nextID++
	return cloneTodo(todo), nil
}

func (r *todoRepo) UpdateTodo(_ context.Context, todo *biz.Todo) (*biz.Todo, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	current, ok := r.todos[todo.ID]
	if !ok {
		return nil, biz.ErrTodoNotFound
	}
	updated := cloneTodo(todo)
	updated.CreateTime = current.CreateTime
	updated.UpdateTime = time.Now()
	r.todos[updated.ID] = cloneTodo(updated)
	return cloneTodo(updated), nil
}

func (r *todoRepo) DeleteTodo(_ context.Context, id int64) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.todos[id]; !ok {
		return biz.ErrTodoNotFound
	}
	delete(r.todos, id)
	return nil
}

func cloneTodo(todo *biz.Todo) *biz.Todo {
	if todo == nil {
		return nil
	}
	return &biz.Todo{
		ID:         todo.ID,
		Title:      todo.Title,
		Content:    todo.Content,
		Completed:  todo.Completed,
		CreateTime: todo.CreateTime,
		UpdateTime: todo.UpdateTime,
	}
}
