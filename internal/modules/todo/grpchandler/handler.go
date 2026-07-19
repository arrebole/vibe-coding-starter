package grpchandler

import (
	"context"

	"github.com/arrebole/vibe-coding-starter/internal/modules/todo/dto"
	"github.com/arrebole/vibe-coding-starter/internal/modules/todo/service"
	todov1 "github.com/arrebole/vibe-coding-starter/proto/public/todo/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

const maxListLimit = 100

type Handler struct {
	todov1.UnimplementedTodoServiceServer
	service *service.Service
}

func New(service *service.Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) ListTodos(ctx context.Context, request *todov1.ListTodosRequest) (*todov1.ListTodosResponse, error) {
	limit := int(request.GetLimit())
	if limit < 0 || limit > maxListLimit {
		return nil, status.Errorf(codes.InvalidArgument, "limit 必须在 0 到 %d 之间", maxListLimit)
	}

	todos, err := h.service.List(ctx, limit)
	if err != nil {
		return nil, err
	}

	return &todov1.ListTodosResponse{Todos: toProtoTodos(todos)}, nil
}

func toProtoTodos(todos []dto.TodoResponse) []*todov1.Todo {
	items := make([]*todov1.Todo, 0, len(todos))
	for _, todo := range todos {
		items = append(items, &todov1.Todo{
			Id:        uint64(todo.ID),
			Title:     todo.Title,
			Completed: todo.Completed,
			CreatedAt: timestamppb.New(todo.CreatedAt),
			UpdatedAt: timestamppb.New(todo.UpdatedAt),
		})
	}
	return items
}
