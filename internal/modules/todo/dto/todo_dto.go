package dto

import (
	"time"

	"github.com/arrebole/vibe-coding-starter/internal/modules/todo/entity"
)

type CreateTodoInput struct {
	Title string
}

type UpdateTodoInput struct {
	ID        uint
	Title     *string
	Completed *bool
}

type TodoResponse struct {
	ID        uint      `json:"id"`
	Title     string    `json:"title"`
	Completed bool      `json:"completed"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func FromEntity(todo entity.Todo) TodoResponse {
	return TodoResponse{
		ID:        todo.ID,
		Title:     todo.Title,
		Completed: todo.Completed,
		CreatedAt: todo.CreatedAt,
		UpdatedAt: todo.UpdatedAt,
	}
}

func FromEntities(todos []entity.Todo) []TodoResponse {
	items := make([]TodoResponse, 0, len(todos))
	for _, todo := range todos {
		items = append(items, FromEntity(todo))
	}
	return items
}
