package service

import (
	"context"
	"log/slog"
	"testing"
	"time"

	"github.com/arrebole/vibe-coding-starter/internal/grpcclient"
	"github.com/arrebole/vibe-coding-starter/internal/modules/todo/dto"
	"github.com/arrebole/vibe-coding-starter/internal/modules/todo/entity"
	"github.com/arrebole/vibe-coding-starter/internal/modules/todo/repository"
)

func TestCreateTodoTrimsTitleAndClearsCache(t *testing.T) {
	repo := &fakeRepository{}
	cache := &fakeTodoCache{}
	svc := New(repo, cache, fakeDemoClient{}, slog.Default())

	todo, err := svc.Create(context.Background(), dto.CreateTodoInput{Title: "  写模板  "})
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}

	if todo.Title != "写模板" {
		t.Fatalf("Title = %q, want %q", todo.Title, "写模板")
	}
	if repo.created == nil {
		t.Fatal("expected repository Create to be called")
	}
	if !cache.deleted {
		t.Fatal("expected list cache to be cleared")
	}
}

func TestCreateTodoRejectsBlankTitle(t *testing.T) {
	svc := New(&fakeRepository{}, nil, nil, slog.Default())

	_, err := svc.Create(context.Background(), dto.CreateTodoInput{Title: "   "})
	if err != ErrInvalidTitle {
		t.Fatalf("Create() error = %v, want %v", err, ErrInvalidTitle)
	}
}

func TestListTodoUsesCache(t *testing.T) {
	expected := []dto.TodoResponse{{ID: 1, Title: "缓存任务"}}
	repo := &fakeRepository{}
	cache := &fakeTodoCache{cached: expected, hit: true}
	svc := New(repo, cache, nil, slog.Default())

	todos, err := svc.List(context.Background(), 10)
	if err != nil {
		t.Fatalf("List() error = %v", err)
	}
	if len(todos) != 1 || todos[0].Title != expected[0].Title {
		t.Fatalf("List() = %#v, want %#v", todos, expected)
	}
	if repo.listCalled {
		t.Fatal("expected repository not to be called when cache hits")
	}
}

type fakeRepository struct {
	created    *entity.Todo
	listCalled bool
}

func (r *fakeRepository) List(context.Context, int) ([]entity.Todo, error) {
	r.listCalled = true
	return []entity.Todo{{ID: 1, Title: "数据库任务", CreatedAt: time.Now(), UpdatedAt: time.Now()}}, nil
}

func (r *fakeRepository) Create(_ context.Context, todo *entity.Todo) error {
	todo.ID = 1
	todo.CreatedAt = time.Now()
	todo.UpdatedAt = todo.CreatedAt
	r.created = todo
	return nil
}

func (r *fakeRepository) FindByID(context.Context, uint) (*entity.Todo, error) {
	return nil, repository.ErrNotFound
}

func (r *fakeRepository) Update(context.Context, *entity.Todo) error {
	return nil
}

func (r *fakeRepository) Delete(context.Context, uint) error {
	return nil
}

type fakeTodoCache struct {
	cached  []dto.TodoResponse
	hit     bool
	deleted bool
}

func (c *fakeTodoCache) GetList(context.Context) ([]dto.TodoResponse, bool, error) {
	return c.cached, c.hit, nil
}

func (c *fakeTodoCache) SetList(context.Context, []dto.TodoResponse) error {
	return nil
}

func (c *fakeTodoCache) DeleteList(context.Context) error {
	c.deleted = true
	return nil
}

type fakeDemoClient struct{}

func (fakeDemoClient) Ping(context.Context, string) (string, error) {
	return "", grpcclient.ErrDemoClientNotConfigured
}

func (fakeDemoClient) Close() error {
	return nil
}
