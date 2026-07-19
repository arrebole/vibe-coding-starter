package grpchandler

import (
	"context"
	"log/slog"
	"net"
	"testing"
	"time"

	"github.com/arrebole/vibe-coding-starter/internal/modules/todo/entity"
	"github.com/arrebole/vibe-coding-starter/internal/modules/todo/service"
	todov1 "github.com/arrebole/vibe-coding-starter/proto/public/todo/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
	"google.golang.org/grpc/test/bufconn"
)

func TestListTodosReturnsTodos(t *testing.T) {
	client, repository, cleanup := startTodoService(t)
	defer cleanup()

	createdAt := time.Date(2026, 7, 19, 10, 0, 0, 0, time.UTC)
	updatedAt := createdAt.Add(time.Hour)
	repository.todos = []entity.Todo{{
		ID:        1,
		Title:     "写 gRPC",
		Completed: true,
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	}}

	response, err := client.ListTodos(context.Background(), &todov1.ListTodosRequest{Limit: 20})
	if err != nil {
		t.Fatalf("ListTodos() error = %v", err)
	}
	if repository.limit != 20 {
		t.Fatalf("repository limit = %d, want %d", repository.limit, 20)
	}
	if len(response.GetTodos()) != 1 {
		t.Fatalf("ListTodos() returned %d todos, want 1", len(response.GetTodos()))
	}

	todo := response.GetTodos()[0]
	if todo.GetId() != 1 || todo.GetTitle() != "写 gRPC" || !todo.GetCompleted() {
		t.Fatalf("ListTodos() todo = %#v", todo)
	}
	if !todo.GetCreatedAt().AsTime().Equal(createdAt) {
		t.Fatalf("created_at = %v, want %v", todo.GetCreatedAt().AsTime(), createdAt)
	}
	if !todo.GetUpdatedAt().AsTime().Equal(updatedAt) {
		t.Fatalf("updated_at = %v, want %v", todo.GetUpdatedAt().AsTime(), updatedAt)
	}
}

func TestListTodosRejectsInvalidLimit(t *testing.T) {
	client, _, cleanup := startTodoService(t)
	defer cleanup()

	_, err := client.ListTodos(context.Background(), &todov1.ListTodosRequest{Limit: 101})
	if status.Code(err) != codes.InvalidArgument {
		t.Fatalf("ListTodos() code = %v, want %v", status.Code(err), codes.InvalidArgument)
	}
}

type fakeRepository struct {
	todos []entity.Todo
	limit int
}

func (r *fakeRepository) List(_ context.Context, limit int) ([]entity.Todo, error) {
	r.limit = limit
	return r.todos, nil
}

func (r *fakeRepository) Create(context.Context, *entity.Todo) error {
	return nil
}

func (r *fakeRepository) FindByID(context.Context, uint) (*entity.Todo, error) {
	return nil, nil
}

func (r *fakeRepository) Update(context.Context, *entity.Todo) error {
	return nil
}

func (r *fakeRepository) Delete(context.Context, uint) error {
	return nil
}

func startTodoService(t *testing.T) (todov1.TodoServiceClient, *fakeRepository, func()) {
	t.Helper()

	repository := &fakeRepository{}
	todoService := service.New(repository, nil, slog.Default())
	listener := bufconn.Listen(1024 * 1024)
	server := grpc.NewServer()
	todov1.RegisterTodoServiceServer(server, New(todoService))

	go func() {
		_ = server.Serve(listener)
	}()

	conn, err := grpc.NewClient(
		"passthrough:///bufnet",
		grpc.WithContextDialer(func(context.Context, string) (net.Conn, error) {
			return listener.Dial()
		}),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		t.Fatalf("grpc.NewClient() error = %v", err)
	}

	return todov1.NewTodoServiceClient(conn), repository, func() {
		_ = conn.Close()
		server.Stop()
		_ = listener.Close()
	}
}
