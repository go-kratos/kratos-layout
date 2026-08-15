package service

import (
	"context"
	"errors"
	"io"
	"strings"

	v1 "github.com/go-kratos/kratos-layout/api/todo/v1"
	"github.com/go-kratos/kratos-layout/internal/biz"

	"github.com/google/uuid"
	"go.einride.tech/aip/fieldmask"
	"go.einride.tech/aip/filtering"
	"go.einride.tech/aip/ordering"
	"go.einride.tech/aip/pagination"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

const (
	defaultPageSize = 20
)

// TodoService is a todo service.
type TodoService struct {
	v1.UnimplementedTodoServiceServer

	uc *biz.TodoUsecase
}

// NewTodoService new a todo service.
func NewTodoService(uc *biz.TodoUsecase) *TodoService {
	return &TodoService{uc: uc}
}

// CreateTodo creates a todo item.
func (s *TodoService) CreateTodo(ctx context.Context, req *v1.CreateTodoRequest) (*v1.Todo, error) {
	todo, err := s.uc.CreateTodo(ctx, convertTodo(req.GetTodo()))
	if err != nil {
		return nil, err
	}
	return convertTodoReply(todo), nil
}

// GetTodo returns a todo item by ID.
func (s *TodoService) GetTodo(ctx context.Context, req *v1.GetTodoRequest) (*v1.Todo, error) {
	id, err := parseTodoID(req.GetId())
	if err != nil {
		return nil, err
	}
	todo, err := s.uc.GetTodo(ctx, id)
	if err != nil {
		return nil, err
	}
	return convertTodoReply(todo), nil
}

// ListTodos lists todo items.
func (s *TodoService) ListTodos(ctx context.Context, req *v1.ListTodosRequest) (*v1.TodoSet, error) {
	declarations, err := filtering.NewDeclarations(
		filtering.DeclareStandardFunctions(),
		filtering.DeclareIdent("id", filtering.TypeString),
		filtering.DeclareIdent("title", filtering.TypeString),
		filtering.DeclareIdent("content", filtering.TypeString),
		filtering.DeclareIdent("completed", filtering.TypeBool),
		filtering.DeclareIdent("created_at", filtering.TypeTimestamp),
		filtering.DeclareIdent("updated_at", filtering.TypeTimestamp),
	)
	if err != nil {
		return nil, err
	}
	filter, err := filtering.ParseFilter(req, declarations)
	if err != nil {
		return nil, err
	}
	pageToken, err := pagination.ParsePageToken(req)
	if err != nil {
		return nil, err
	}
	orderBy, err := ordering.ParseOrderBy(req)
	if err != nil {
		return nil, err
	}
	if err := orderBy.ValidateForPaths("id", "title", "created_at", "updated_at"); err != nil {
		return nil, err
	}
	if req.PageSize <= 0 {
		req.PageSize = defaultPageSize
	}
	todos, err := s.uc.ListTodos(ctx,
		biz.ListFilter(filter),
		biz.ListOrderBy(orderBy),
		biz.ListLimit(int(req.PageSize)),
		biz.ListOffset(int(pageToken.Offset)),
	)
	if err != nil {
		return nil, err
	}
	set := &v1.TodoSet{
		Todos: make([]*v1.Todo, 0, len(todos)),
	}
	if len(todos) >= int(req.PageSize) {
		set.NextPageToken = pageToken.Next(req).String()
	}
	for _, todo := range todos {
		set.Todos = append(set.Todos, convertTodoReply(todo))
	}
	return set, nil
}

// UpdateTodo updates a todo item.
func (s *TodoService) UpdateTodo(ctx context.Context, req *v1.UpdateTodoRequest) (*v1.Todo, error) {
	if req.GetTodo().GetId() == "" || req.GetUpdateMask() == nil || len(req.GetUpdateMask().GetPaths()) == 0 {
		return nil, biz.ErrTodoInvalidArgument
	}
	current, err := s.GetTodo(ctx, &v1.GetTodoRequest{Id: req.GetTodo().GetId()})
	if err != nil {
		return nil, err
	}
	fieldmask.Update(req.GetUpdateMask(), current, req.GetTodo())
	todo, err := s.uc.UpdateTodo(ctx, convertTodo(current))
	if err != nil {
		return nil, err
	}
	return convertTodoReply(todo), nil
}

// DeleteTodo soft-deletes a todo item.
func (s *TodoService) DeleteTodo(ctx context.Context, req *v1.DeleteTodoRequest) (*emptypb.Empty, error) {
	id, err := parseTodoID(req.GetId())
	if err != nil {
		return nil, err
	}
	if err := s.uc.DeleteTodo(ctx, id); err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

// WatchTodos streams todo snapshots from the server to the client.
func (s *TodoService) WatchTodos(req *v1.WatchTodosRequest, stream v1.TodoService_WatchTodosServer) error {
	declarations, err := filtering.NewDeclarations(
		filtering.DeclareStandardFunctions(),
		filtering.DeclareIdent("id", filtering.TypeString),
		filtering.DeclareIdent("title", filtering.TypeString),
		filtering.DeclareIdent("content", filtering.TypeString),
		filtering.DeclareIdent("completed", filtering.TypeBool),
		filtering.DeclareIdent("created_at", filtering.TypeTimestamp),
		filtering.DeclareIdent("updated_at", filtering.TypeTimestamp),
	)
	if err != nil {
		return err
	}
	filter, err := filtering.ParseFilter(req, declarations)
	if err != nil {
		return err
	}
	pageToken, err := pagination.ParsePageToken(req)
	if err != nil {
		return err
	}
	orderBy, err := ordering.ParseOrderBy(req)
	if err != nil {
		return err
	}
	if err := orderBy.ValidateForPaths("id", "title", "created_at", "updated_at"); err != nil {
		return err
	}
	if req.PageSize <= 0 {
		req.PageSize = defaultPageSize
	}
	todos, err := s.uc.ListTodos(stream.Context(),
		biz.ListFilter(filter),
		biz.ListOrderBy(orderBy),
		biz.ListLimit(int(req.PageSize)),
		biz.ListOffset(int(pageToken.Offset)),
	)
	if err != nil {
		return err
	}
	for _, todo := range todos {
		if err := stream.Send(newTodoEvent("snapshot", convertTodoReply(todo))); err != nil {
			return err
		}
	}
	return nil
}

// SyncTodos exchanges todo changes in both directions.
func (s *TodoService) SyncTodos(stream v1.TodoService_SyncTodosServer) error {
	for {
		req, err := stream.Recv()
		if errors.Is(err, io.EOF) {
			return nil
		}
		if err != nil {
			return err
		}
		var event *v1.TodoEvent
		switch strings.ToLower(req.GetAction()) {
		case "create":
			todo, err := s.CreateTodo(stream.Context(), &v1.CreateTodoRequest{Todo: req.GetTodo()})
			if err != nil {
				return err
			}
			event = newTodoEvent("created", todo)
		case "update":
			todo, err := s.UpdateTodo(stream.Context(), &v1.UpdateTodoRequest{
				Todo:       req.GetTodo(),
				UpdateMask: req.GetUpdateMask(),
			})
			if err != nil {
				return err
			}
			event = newTodoEvent("updated", todo)
		case "delete":
			id := req.GetId()
			if id == "" {
				id = req.GetTodo().GetId()
			}
			if _, err := s.DeleteTodo(stream.Context(), &v1.DeleteTodoRequest{Id: id}); err != nil {
				return err
			}
			// The record is hidden from reads now, so only the id and the
			// resulting status are meaningful.
			event = newTodoEvent("deleted", &v1.Todo{
				Id:     id,
				Status: v1.TodoStatus_TODO_STATUS_DELETED,
			})
		default:
			return biz.ErrTodoInvalidArgument
		}
		if err := stream.Send(event); err != nil {
			return err
		}
	}
}

// parseTodoID validates an incoming todo id at the service boundary.
func parseTodoID(id string) (uuid.UUID, error) {
	parsed, err := uuid.Parse(id)
	if err != nil {
		return uuid.Nil, biz.ErrTodoInvalidArgument
	}
	return parsed, nil
}

// convertTodo parses an incoming proto into a DO. Timestamps are omitted: they
// are server-assigned and read-only, so only the mutable fields cross over.
//
// The id is left at uuid.Nil when absent or unparseable. Callers that require a
// real id validate it with parseTodoID first; on the create path no id is
// supplied and the storage layer assigns one.
func convertTodo(in *v1.Todo) *biz.Todo {
	if in == nil {
		return nil
	}
	id, _ := uuid.Parse(in.GetId())
	return &biz.Todo{
		ID:        id,
		Title:     in.GetTitle(),
		Content:   in.GetContent(),
		Completed: in.GetCompleted(),
	}
}

// newTodoEvent wraps an outgoing todo in a stream event. It takes the reply
// shape directly so handlers that already hold one do not round-trip it back
// through a DO, which would drop the server-assigned timestamps.
func newTodoEvent(action string, todo *v1.Todo) *v1.TodoEvent {
	return &v1.TodoEvent{
		Action:    action,
		Todo:      todo,
		EventTime: timestamppb.Now(),
	}
}

func convertTodoReply(in *biz.Todo) *v1.Todo {
	if in == nil {
		return nil
	}
	return &v1.Todo{
		Id:        in.ID.String(),
		Title:     in.Title,
		Content:   in.Content,
		Completed: in.Completed,
		Status:    convertTodoStatus(in.Status),
		CreatedAt: timestamppb.New(in.CreatedAt),
		UpdatedAt: timestamppb.New(in.UpdatedAt),
	}
}

// convertTodoStatus maps a domain status onto the api enum. The two share their
// numeric values, but the mapping is written out so neither side can drift into
// the other silently.
func convertTodoStatus(in biz.TodoStatus) v1.TodoStatus {
	switch in {
	case biz.TodoStatusActive:
		return v1.TodoStatus_TODO_STATUS_ACTIVE
	case biz.TodoStatusDeleted:
		return v1.TodoStatus_TODO_STATUS_DELETED
	default:
		return v1.TodoStatus_TODO_STATUS_UNSPECIFIED
	}
}
