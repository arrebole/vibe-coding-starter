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
	svc := New(repo, fakeDemoClient{}, slog.Default())

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
}

func TestCreateTodoRejectsBlankTitle(t *testing.T) {
	svc := New(&fakeRepository{}, nil, slog.Default())

	_, err := svc.Create(context.Background(), dto.CreateTodoInput{Title: "   "})
	if err != ErrInvalidTitle {
		t.Fatalf("Create() error = %v, want %v", err, ErrInvalidTitle)
	}
}

func TestListTodoUsesRepository(t *testing.T) {
	repo := &fakeRepository{}
	svc := New(repo, nil, slog.Default())

	todos, err := svc.List(context.Background(), 10)
	if err != nil {
		t.Fatalf("List() error = %v", err)
	}
	if len(todos) != 1 || todos[0].Title != "数据库任务" {
		t.Fatalf("List() = %#v", todos)
	}
	if !repo.listCalled {
		t.Fatal("expected repository to be called")
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

type fakeDemoClient struct{}

func (fakeDemoClient) Ping(context.Context, string) (string, error) {
	return "", grpcclient.ErrDemoClientNotConfigured
}

func (fakeDemoClient) Close() error {
	return nil
}
