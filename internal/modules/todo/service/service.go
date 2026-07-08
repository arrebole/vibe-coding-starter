package service

import (
	"context"
	"errors"
	"log/slog"
	"strings"

	"github.com/arrebole/vibe-coding-starter/internal/grpcclient"
	"github.com/arrebole/vibe-coding-starter/internal/modules/todo/dto"
	"github.com/arrebole/vibe-coding-starter/internal/modules/todo/entity"
	"github.com/arrebole/vibe-coding-starter/internal/modules/todo/repository"
)

var (
	ErrInvalidTitle = errors.New("todo 标题不能为空")
	ErrNotFound     = errors.New("todo 不存在")
)

type Service struct {
	repository repository.Repository
	cache      TodoCache
	demoClient grpcclient.DemoClient
	log        *slog.Logger
}

type TodoCache interface {
	GetList(ctx context.Context) ([]dto.TodoResponse, bool, error)
	SetList(ctx context.Context, todos []dto.TodoResponse) error
	DeleteList(ctx context.Context) error
}

func New(repository repository.Repository, cache TodoCache, demoClient grpcclient.DemoClient, log *slog.Logger) *Service {
	return &Service{
		repository: repository,
		cache:      cache,
		demoClient: demoClient,
		log:        log,
	}
}

func (s *Service) List(ctx context.Context, limit int) ([]dto.TodoResponse, error) {
	if s.cache != nil {
		if todos, ok, err := s.cache.GetList(ctx); err == nil && ok {
			return todos, nil
		} else if err != nil {
			s.log.Warn("读取 Todo 缓存失败", slog.Any("error", err))
		}
	}

	todos, err := s.repository.List(ctx, limit)
	if err != nil {
		return nil, err
	}

	output := dto.FromEntities(todos)
	if s.cache != nil {
		if err := s.cache.SetList(ctx, output); err != nil {
			s.log.Warn("写入 Todo 缓存失败", slog.Any("error", err))
		}
	}
	return output, nil
}

func (s *Service) Create(ctx context.Context, input dto.CreateTodoInput) (dto.TodoResponse, error) {
	title := strings.TrimSpace(input.Title)
	if title == "" {
		return dto.TodoResponse{}, ErrInvalidTitle
	}

	todo := &entity.Todo{Title: title}
	if err := s.repository.Create(ctx, todo); err != nil {
		return dto.TodoResponse{}, err
	}
	s.clearListCache(ctx)
	s.callExternalDemo(ctx, title)
	return dto.FromEntity(*todo), nil
}

func (s *Service) Update(ctx context.Context, input dto.UpdateTodoInput) (dto.TodoResponse, error) {
	todo, err := s.repository.FindByID(ctx, input.ID)
	if errors.Is(err, repository.ErrNotFound) {
		return dto.TodoResponse{}, ErrNotFound
	}
	if err != nil {
		return dto.TodoResponse{}, err
	}

	if input.Title != nil {
		title := strings.TrimSpace(*input.Title)
		if title == "" {
			return dto.TodoResponse{}, ErrInvalidTitle
		}
		todo.Title = title
	}
	if input.Completed != nil {
		todo.Completed = *input.Completed
	}

	if err := s.repository.Update(ctx, todo); err != nil {
		return dto.TodoResponse{}, err
	}
	s.clearListCache(ctx)
	return dto.FromEntity(*todo), nil
}

func (s *Service) Delete(ctx context.Context, id uint) error {
	err := s.repository.Delete(ctx, id)
	if errors.Is(err, repository.ErrNotFound) {
		return ErrNotFound
	}
	if err != nil {
		return err
	}
	s.clearListCache(ctx)
	return nil
}

func (s *Service) clearListCache(ctx context.Context) {
	if s.cache == nil {
		return
	}
	if err := s.cache.DeleteList(ctx); err != nil {
		s.log.Warn("清理 Todo 缓存失败", slog.Any("error", err))
	}
}

func (s *Service) callExternalDemo(ctx context.Context, title string) {
	if s.demoClient == nil {
		return
	}
	if _, err := s.demoClient.Ping(ctx, title); err != nil && !errors.Is(err, grpcclient.ErrDemoClientNotConfigured) {
		s.log.Warn("调用外部 demo gRPC 服务失败", slog.Any("error", err))
	}
}
