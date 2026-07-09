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
	demoClient grpcclient.DemoClient
	log        *slog.Logger
}

func New(repository repository.Repository, demoClient grpcclient.DemoClient, log *slog.Logger) *Service {
	return &Service{
		repository: repository,
		demoClient: demoClient,
		log:        log,
	}
}

func (s *Service) List(ctx context.Context, limit int) ([]dto.TodoResponse, error) {
	todos, err := s.repository.List(ctx, limit)
	if err != nil {
		return nil, err
	}
	return dto.FromEntities(todos), nil
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
	return nil
}

func (s *Service) callExternalDemo(ctx context.Context, title string) {
	if s.demoClient == nil {
		return
	}
	if _, err := s.demoClient.Ping(ctx, title); err != nil && !errors.Is(err, grpcclient.ErrDemoClientNotConfigured) {
		s.log.Warn("调用外部 demo gRPC 服务失败", slog.Any("error", err))
	}
}
