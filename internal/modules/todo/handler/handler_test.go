package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"

	"github.com/arrebole/vibe-coding-starter/internal/modules/todo/dto"
	"github.com/arrebole/vibe-coding-starter/internal/modules/todo/entity"
	"github.com/arrebole/vibe-coding-starter/internal/modules/todo/repository"
	"github.com/arrebole/vibe-coding-starter/internal/modules/todo/service"
)

func TestCreateTodo(t *testing.T) {
	gin.SetMode(gin.TestMode)
	engine := gin.New()
	h := New(service.New(&handlerFakeRepository{}, nil, slog.Default()))
	group := engine.Group("/api/v1")
	h.Register(group)

	body := bytes.NewBufferString(`{"title":"写测试"}`)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/todos", body)
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	engine.ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("status = %d, body = %s", rec.Code, rec.Body.String())
	}

	var payload struct {
		Code int              `json:"code"`
		Data dto.TodoResponse `json:"data"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &payload); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if payload.Code != 0 || payload.Data.Title != "写测试" {
		t.Fatalf("payload = %#v", payload)
	}
}

func TestCreateTodoRejectsInvalidBody(t *testing.T) {
	gin.SetMode(gin.TestMode)
	engine := gin.New()
	h := New(service.New(&handlerFakeRepository{}, nil, slog.Default()))
	group := engine.Group("/api/v1")
	h.Register(group)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/todos", bytes.NewBufferString(`{}`))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	engine.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusBadRequest)
	}
}

type handlerFakeRepository struct{}

func (handlerFakeRepository) List(context.Context, int) ([]entity.Todo, error) {
	return nil, nil
}

func (handlerFakeRepository) Create(_ context.Context, todo *entity.Todo) error {
	todo.ID = 1
	return nil
}

func (handlerFakeRepository) FindByID(context.Context, uint) (*entity.Todo, error) {
	return nil, repository.ErrNotFound
}

func (handlerFakeRepository) Update(context.Context, *entity.Todo) error {
	return nil
}

func (handlerFakeRepository) Delete(context.Context, uint) error {
	return nil
}
