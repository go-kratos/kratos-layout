package data

import (
	"context"
	stdsql "database/sql"
	"testing"

	"github.com/go-kratos/kratos-layout/internal/biz"
	"github.com/go-kratos/kratos-layout/internal/data/ent"

	"entgo.io/ent/dialect"
	entsql "entgo.io/ent/dialect/sql"
	kratoserrors "github.com/go-kratos/kratos/v3/errors"
	"github.com/google/uuid"
	"go.einride.tech/aip/ordering"
	_ "modernc.org/sqlite"
)

// newTestRepo opens an in-memory SQLite database with the schema applied, so
// repo tests exercise real ent queries at the storage boundary. The modernc
// driver registers itself as "sqlite", so the ent dialect is set explicitly.
func newTestRepo(t *testing.T) (biz.TodoRepo, *ent.Client) {
	t.Helper()
	db, err := stdsql.Open("sqlite", "file:"+t.Name()+"?mode=memory&cache=shared&_pragma=foreign_keys(1)")
	if err != nil {
		t.Fatalf("open sqlite: %v", err)
	}
	client := ent.NewClient(ent.Driver(entsql.OpenDB(dialect.SQLite, db)))
	t.Cleanup(func() {
		client.Close()
	})
	if err := client.Schema.Create(context.Background()); err != nil {
		t.Fatalf("create schema: %v", err)
	}
	return NewTodoRepo(&Data{db: client}), client
}

// orderByCreatedAt keeps list assertions deterministic.
func orderByCreatedAt() biz.ListOption {
	return biz.ListOrderBy(ordering.OrderBy{
		Fields: []ordering.Field{{Path: "created_at"}},
	})
}

func TestTodoRepoCRUD(t *testing.T) {
	ctx := context.Background()
	repo, _ := newTestRepo(t)

	created, err := repo.CreateTodo(ctx, &biz.Todo{
		Title:   "write tests",
		Content: "cover repo CRUD",
	})
	if err != nil {
		t.Fatalf("CreateTodo() error = %v", err)
	}
	if created.ID == uuid.Nil {
		t.Fatal("CreateTodo() did not assign an id")
	}
	if created.CreatedAt.IsZero() || created.UpdatedAt.IsZero() {
		t.Fatal("CreateTodo() did not set timestamps")
	}

	got, err := repo.FindByID(ctx, created.ID)
	if err != nil {
		t.Fatalf("FindByID() error = %v", err)
	}
	if got.Title != "write tests" || got.Content != "cover repo CRUD" {
		t.Fatalf("FindByID() = %+v, want created todo", got)
	}

	updated, err := repo.UpdateTodo(ctx, &biz.Todo{
		ID:        created.ID,
		Title:     "write repo tests",
		Content:   "cover repo CRUD",
		Completed: true,
	})
	if err != nil {
		t.Fatalf("UpdateTodo() error = %v", err)
	}
	if updated.Title != "write repo tests" || !updated.Completed {
		t.Fatalf("UpdateTodo() = %+v, want updated title and completed", updated)
	}
	if updated.ID != created.ID {
		t.Fatalf("UpdateTodo() id = %v, want %v", updated.ID, created.ID)
	}

	if _, err := repo.FindByID(ctx, uuid.Must(uuid.NewV7())); !kratoserrors.IsNotFound(err) {
		t.Fatalf("FindByID(missing) error = %v, want not found", err)
	}
}

func TestTodoRepoDeleteIsSoft(t *testing.T) {
	ctx := context.Background()
	repo, client := newTestRepo(t)

	created, err := repo.CreateTodo(ctx, &biz.Todo{Title: "soft delete me"})
	if err != nil {
		t.Fatalf("CreateTodo() error = %v", err)
	}

	if err := repo.DeleteTodo(ctx, created.ID); err != nil {
		t.Fatalf("DeleteTodo() error = %v", err)
	}

	// The row survives with status flipped to deleted.
	po, err := client.Todo.Get(ctx, created.ID)
	if err != nil {
		t.Fatalf("ent Get() after delete error = %v", err)
	}
	if po.Status != biz.TodoStatusDeleted {
		t.Fatalf("status after delete = %d, want %d", po.Status, biz.TodoStatusDeleted)
	}

	// Every read and write path hides it.
	if _, err := repo.FindByID(ctx, created.ID); !kratoserrors.IsNotFound(err) {
		t.Fatalf("FindByID() after delete error = %v, want not found", err)
	}
	todos, err := repo.ListTodos(ctx, biz.ListLimit(10))
	if err != nil {
		t.Fatalf("ListTodos() error = %v", err)
	}
	if len(todos) != 0 {
		t.Fatalf("ListTodos() len = %d, want 0", len(todos))
	}
	if _, err := repo.UpdateTodo(ctx, &biz.Todo{ID: created.ID, Title: "revive"}); !kratoserrors.IsNotFound(err) {
		t.Fatalf("UpdateTodo() after delete error = %v, want not found", err)
	}
	if err := repo.DeleteTodo(ctx, created.ID); !kratoserrors.IsNotFound(err) {
		t.Fatalf("DeleteTodo() twice error = %v, want not found", err)
	}
}

func TestTodoRepoListPagination(t *testing.T) {
	ctx := context.Background()
	repo, _ := newTestRepo(t)

	for _, title := range []string{"first", "second", "third"} {
		if _, err := repo.CreateTodo(ctx, &biz.Todo{Title: title}); err != nil {
			t.Fatalf("CreateTodo(%q) error = %v", title, err)
		}
	}

	firstPage, err := repo.ListTodos(ctx, orderByCreatedAt(), biz.ListLimit(2))
	if err != nil {
		t.Fatalf("ListTodos(first page) error = %v", err)
	}
	if len(firstPage) != 2 {
		t.Fatalf("ListTodos(first page) len = %d, want 2", len(firstPage))
	}
	if firstPage[0].Title != "first" {
		t.Fatalf("ListTodos(first page) title = %q, want first", firstPage[0].Title)
	}

	secondPage, err := repo.ListTodos(ctx, orderByCreatedAt(), biz.ListLimit(2), biz.ListOffset(2))
	if err != nil {
		t.Fatalf("ListTodos(second page) error = %v", err)
	}
	if len(secondPage) != 1 {
		t.Fatalf("ListTodos(second page) len = %d, want 1", len(secondPage))
	}
	if secondPage[0].Title != "third" {
		t.Fatalf("ListTodos(second page) title = %q, want third", secondPage[0].Title)
	}

	if _, err := repo.ListTodos(ctx, biz.ListLimit(0)); !kratoserrors.IsBadRequest(err) {
		t.Fatalf("ListTodos(zero limit) error = %v, want bad request", err)
	}
}

// Offset pagination is only correct under a total order. With no order_by the
// repo falls back to id, which is UUIDv7 and therefore time-ordered, so paging
// covers every row exactly once.
func TestTodoRepoListDefaultOrderIsStable(t *testing.T) {
	ctx := context.Background()
	repo, _ := newTestRepo(t)

	const total = 6
	titles := make([]string, 0, total)
	for i := range total {
		title := string(rune('a' + i))
		if _, err := repo.CreateTodo(ctx, &biz.Todo{Title: title}); err != nil {
			t.Fatalf("CreateTodo(%q) error = %v", title, err)
		}
		titles = append(titles, title)
	}

	// Page through without requesting any ordering.
	var paged []string
	for offset := 0; offset < total; offset += 2 {
		page, err := repo.ListTodos(ctx, biz.ListLimit(2), biz.ListOffset(offset))
		if err != nil {
			t.Fatalf("ListTodos(offset=%d) error = %v", offset, err)
		}
		for _, todo := range page {
			paged = append(paged, todo.Title)
		}
	}

	if len(paged) != total {
		t.Fatalf("paged through %d rows, want %d (duplicated or skipped)", len(paged), total)
	}
	for i, want := range titles {
		if paged[i] != want {
			t.Fatalf("paged = %v, want creation order %v", paged, titles)
		}
	}
}

// A caller-supplied sort key wins, and id breaks ties so equal keys still page
// deterministically.
func TestTodoRepoListOrderByTiesBrokenByID(t *testing.T) {
	ctx := context.Background()
	repo, _ := newTestRepo(t)

	// Every row shares one title, so the requested sort key cannot order them.
	var ids []uuid.UUID
	for range 4 {
		created, err := repo.CreateTodo(ctx, &biz.Todo{Title: "same"})
		if err != nil {
			t.Fatalf("CreateTodo() error = %v", err)
		}
		ids = append(ids, created.ID)
	}

	got, err := repo.ListTodos(ctx,
		biz.ListOrderBy(ordering.OrderBy{Fields: []ordering.Field{{Path: "title"}}}),
		biz.ListLimit(10),
	)
	if err != nil {
		t.Fatalf("ListTodos() error = %v", err)
	}
	if len(got) != len(ids) {
		t.Fatalf("ListTodos() len = %d, want %d", len(got), len(ids))
	}
	for i, want := range ids {
		if got[i].ID != want {
			t.Fatalf("row %d id = %v, want %v (ties must fall back to id)", i, got[i].ID, want)
		}
	}
}

func TestTodoRepoListOrderByDesc(t *testing.T) {
	ctx := context.Background()
	repo, _ := newTestRepo(t)

	for _, title := range []string{"alpha", "beta", "gamma"} {
		if _, err := repo.CreateTodo(ctx, &biz.Todo{Title: title}); err != nil {
			t.Fatalf("CreateTodo(%q) error = %v", title, err)
		}
	}

	todos, err := repo.ListTodos(ctx,
		biz.ListOrderBy(ordering.OrderBy{
			Fields: []ordering.Field{{Path: "title", Desc: true}},
		}),
		biz.ListLimit(10),
	)
	if err != nil {
		t.Fatalf("ListTodos(order desc) error = %v", err)
	}
	if len(todos) != 3 {
		t.Fatalf("ListTodos(order desc) len = %d, want 3", len(todos))
	}
	if todos[0].Title != "gamma" || todos[2].Title != "alpha" {
		t.Fatalf("ListTodos(order desc) = %q..%q, want gamma..alpha", todos[0].Title, todos[2].Title)
	}
}
