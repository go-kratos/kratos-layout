package service

import (
	"context"
	"io"
	"testing"

	v1 "github.com/go-kratos/kratos-layout/api/todo/v1"
	"github.com/go-kratos/kratos-layout/internal/biz"
	"github.com/go-kratos/kratos-layout/internal/data"

	kratoserrors "github.com/go-kratos/kratos/v3/errors"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/types/known/fieldmaskpb"
)

func newTestTodoService() *TodoService {
	repo := data.NewTodoRepo(&data.Data{})
	uc := biz.NewTodoUsecase(repo)
	return NewTodoService(uc)
}

func TestTodoServiceCRUD(t *testing.T) {
	ctx := context.Background()
	svc := newTestTodoService()

	created, err := svc.CreateTodo(ctx, &v1.CreateTodoRequest{
		Todo: &v1.Todo{
			Title:     "write tests",
			Content:   "cover todo CRUD",
			Completed: false,
		},
	})
	if err != nil {
		t.Fatalf("CreateTodo() error = %v", err)
	}
	if created.GetId() != 1 {
		t.Fatalf("CreateTodo() id = %d, want 1", created.GetId())
	}
	if created.GetCreateTime() == nil || created.GetUpdateTime() == nil {
		t.Fatal("CreateTodo() did not set timestamps")
	}

	got, err := svc.GetTodo(ctx, &v1.GetTodoRequest{Id: created.GetId()})
	if err != nil {
		t.Fatalf("GetTodo() error = %v", err)
	}
	if got.GetTitle() != "write tests" || got.GetContent() != "cover todo CRUD" {
		t.Fatalf("GetTodo() = %+v, want created todo", got)
	}

	updated, err := svc.UpdateTodo(ctx, &v1.UpdateTodoRequest{
		Todo: &v1.Todo{
			Id:        created.GetId(),
			Title:     "write service tests",
			Completed: true,
		},
		UpdateMask: &fieldmaskpb.FieldMask{Paths: []string{"title", "completed"}},
	})
	if err != nil {
		t.Fatalf("UpdateTodo() error = %v", err)
	}
	if updated.GetTitle() != "write service tests" || !updated.GetCompleted() {
		t.Fatalf("UpdateTodo() = %+v, want updated title and completed", updated)
	}
	if updated.GetContent() != "cover todo CRUD" {
		t.Fatalf("UpdateTodo() content = %q, want original content", updated.GetContent())
	}

	if _, err := svc.DeleteTodo(ctx, &v1.DeleteTodoRequest{Id: created.GetId()}); err != nil {
		t.Fatalf("DeleteTodo() error = %v", err)
	}
	if _, err := svc.GetTodo(ctx, &v1.GetTodoRequest{Id: created.GetId()}); !kratoserrors.IsNotFound(err) {
		t.Fatalf("GetTodo() after delete error = %v, want not found", err)
	}
}

func TestTodoServiceListTodosPagination(t *testing.T) {
	ctx := context.Background()
	svc := newTestTodoService()

	for _, title := range []string{"first", "second", "third"} {
		if _, err := svc.CreateTodo(ctx, &v1.CreateTodoRequest{Todo: &v1.Todo{Title: title}}); err != nil {
			t.Fatalf("CreateTodo(%q) error = %v", title, err)
		}
	}

	firstPage, err := svc.ListTodos(ctx, &v1.ListTodosRequest{PageSize: 2})
	if err != nil {
		t.Fatalf("ListTodos(first page) error = %v", err)
	}
	if len(firstPage.GetTodos()) != 2 {
		t.Fatalf("ListTodos(first page) len = %d, want 2", len(firstPage.GetTodos()))
	}
	if firstPage.GetNextPageToken() == "" {
		t.Fatal("ListTodos(first page) next_page_token is empty")
	}

	secondPage, err := svc.ListTodos(ctx, &v1.ListTodosRequest{
		PageSize:  2,
		PageToken: firstPage.GetNextPageToken(),
	})
	if err != nil {
		t.Fatalf("ListTodos(second page) error = %v", err)
	}
	if len(secondPage.GetTodos()) != 1 {
		t.Fatalf("ListTodos(second page) len = %d, want 1", len(secondPage.GetTodos()))
	}
	if secondPage.GetNextPageToken() != "" {
		t.Fatalf("ListTodos(second page) next_page_token = %q, want empty", secondPage.GetNextPageToken())
	}
	if secondPage.GetTodos()[0].GetTitle() != "third" {
		t.Fatalf("ListTodos(second page) title = %q, want third", secondPage.GetTodos()[0].GetTitle())
	}
}

func TestTodoServiceListTodosFilterAndOrderByValidation(t *testing.T) {
	ctx := context.Background()
	svc := newTestTodoService()

	for _, todo := range []*v1.Todo{
		{Title: "write docs", Content: "public docs", Completed: true},
		{Title: "fix bug", Content: "private bug", Completed: false},
		{Title: "fix api", Content: "public api", Completed: true},
	} {
		if _, err := svc.CreateTodo(ctx, &v1.CreateTodoRequest{Todo: todo}); err != nil {
			t.Fatalf("CreateTodo(%q) error = %v", todo.GetTitle(), err)
		}
	}

	reply, err := svc.ListTodos(ctx, &v1.ListTodosRequest{
		PageSize: 10,
		Filter:   `title:"fix" AND completed`,
		OrderBy:  "title desc",
	})
	if err != nil {
		t.Fatalf("ListTodos(filter/order) error = %v", err)
	}
	if len(reply.GetTodos()) != 3 {
		t.Fatalf("ListTodos(filter/order) len = %d, want 3", len(reply.GetTodos()))
	}
	if reply.GetTodos()[0].GetTitle() != "write docs" {
		t.Fatalf("ListTodos(filter/order) first title = %q, want ID order", reply.GetTodos()[0].GetTitle())
	}
}

func TestTodoServiceValidation(t *testing.T) {
	ctx := context.Background()
	svc := newTestTodoService()

	if _, err := svc.CreateTodo(ctx, &v1.CreateTodoRequest{Todo: &v1.Todo{Title: " "}}); !kratoserrors.IsBadRequest(err) {
		t.Fatalf("CreateTodo(empty title) error = %v, want bad request", err)
	}
	if _, err := svc.UpdateTodo(ctx, &v1.UpdateTodoRequest{
		Todo:       &v1.Todo{Id: 1, Title: "missing mask"},
		UpdateMask: &fieldmaskpb.FieldMask{},
	}); !kratoserrors.IsBadRequest(err) {
		t.Fatalf("UpdateTodo(empty mask) error = %v, want bad request", err)
	}
	if _, err := svc.ListTodos(ctx, &v1.ListTodosRequest{PageToken: "bad-token"}); err == nil {
		t.Fatal("ListTodos(bad token) error = nil, want error")
	}
	if _, err := svc.ListTodos(ctx, &v1.ListTodosRequest{Filter: `unknown:"value"`}); err == nil {
		t.Fatal("ListTodos(unsupported filter) error = nil, want error")
	}
	if _, err := svc.ListTodos(ctx, &v1.ListTodosRequest{OrderBy: "content"}); err == nil {
		t.Fatal("ListTodos(unsupported order_by) error = nil, want error")
	}
	if _, err := svc.DeleteTodo(ctx, &v1.DeleteTodoRequest{Id: 1}); !kratoserrors.IsNotFound(err) {
		t.Fatalf("DeleteTodo(missing id) error = %v, want not found", err)
	}
}

func TestTodoServiceWatchTodos(t *testing.T) {
	ctx := context.Background()
	svc := newTestTodoService()

	for _, todo := range []*v1.Todo{
		{Title: "open task", Completed: false},
		{Title: "done task", Completed: true},
	} {
		if _, err := svc.CreateTodo(ctx, &v1.CreateTodoRequest{Todo: todo}); err != nil {
			t.Fatalf("CreateTodo(%q) error = %v", todo.GetTitle(), err)
		}
	}

	stream := &watchTodosStream{fakeServerStream: fakeServerStream{ctx: ctx}}
	if err := svc.WatchTodos(&v1.WatchTodosRequest{
		PageSize: 10,
	}, stream); err != nil {
		t.Fatalf("WatchTodos() error = %v", err)
	}
	if len(stream.events) != 2 {
		t.Fatalf("WatchTodos() events len = %d, want 2", len(stream.events))
	}
	if stream.events[0].GetAction() != "snapshot" || stream.events[0].GetTodo().GetTitle() != "open task" {
		t.Fatalf("WatchTodos() event = %+v, want open snapshot", stream.events[0])
	}
}

func TestTodoServiceSyncTodos(t *testing.T) {
	ctx := context.Background()
	svc := newTestTodoService()
	stream := &syncTodosStream{
		fakeServerStream: fakeServerStream{ctx: ctx},
		requests: []*v1.SyncTodoRequest{
			{
				Action: "create",
				Todo:   &v1.Todo{Title: "streamed todo", Content: "from bidi stream"},
			},
			{
				Action: "update",
				Todo:   &v1.Todo{Id: 1, Completed: true},
				UpdateMask: &fieldmaskpb.FieldMask{
					Paths: []string{"completed"},
				},
			},
			{
				Action: "delete",
				Id:     1,
			},
		},
	}

	if err := svc.SyncTodos(stream); err != nil {
		t.Fatalf("SyncTodos() error = %v", err)
	}
	if got := len(stream.events); got != 3 {
		t.Fatalf("SyncTodos() events len = %d, want 3", got)
	}
	if stream.events[0].GetAction() != "created" || stream.events[0].GetTodo().GetId() != 1 {
		t.Fatalf("SyncTodos() create event = %+v, want created id 1", stream.events[0])
	}
	if stream.events[1].GetAction() != "updated" || !stream.events[1].GetTodo().GetCompleted() {
		t.Fatalf("SyncTodos() update event = %+v, want completed update", stream.events[1])
	}
	if stream.events[2].GetAction() != "deleted" || stream.events[2].GetTodo().GetId() != 1 {
		t.Fatalf("SyncTodos() delete event = %+v, want deleted id 1", stream.events[2])
	}
}

type fakeServerStream struct {
	ctx context.Context
}

func (s fakeServerStream) SetHeader(metadata.MD) error  { return nil }
func (s fakeServerStream) SendHeader(metadata.MD) error { return nil }
func (s fakeServerStream) SetTrailer(metadata.MD)       {}
func (s fakeServerStream) Context() context.Context     { return s.ctx }
func (s fakeServerStream) SendMsg(any) error            { return nil }
func (s fakeServerStream) RecvMsg(any) error            { return nil }

type watchTodosStream struct {
	fakeServerStream
	events []*v1.TodoEvent
}

func (s *watchTodosStream) Send(event *v1.TodoEvent) error {
	s.events = append(s.events, event)
	return nil
}

type syncTodosStream struct {
	fakeServerStream
	requests []*v1.SyncTodoRequest
	events   []*v1.TodoEvent
}

func (s *syncTodosStream) Recv() (*v1.SyncTodoRequest, error) {
	if len(s.requests) == 0 {
		return nil, io.EOF
	}
	req := s.requests[0]
	s.requests = s.requests[1:]
	return req, nil
}

func (s *syncTodosStream) Send(event *v1.TodoEvent) error {
	s.events = append(s.events, event)
	return nil
}
