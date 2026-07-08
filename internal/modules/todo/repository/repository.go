package repository

import (
	"context"
	"errors"

	"gorm.io/gorm"

	"github.com/arrebole/vibe-coding-starter/internal/modules/todo/entity"
)

var ErrNotFound = errors.New("todo 不存在")

type Repository interface {
	List(ctx context.Context, limit int) ([]entity.Todo, error)
	Create(ctx context.Context, todo *entity.Todo) error
	FindByID(ctx context.Context, id uint) (*entity.Todo, error)
	Update(ctx context.Context, todo *entity.Todo) error
	Delete(ctx context.Context, id uint) error
}

type GormRepository struct {
	db *gorm.DB
}

func NewGormRepository(db *gorm.DB) *GormRepository {
	return &GormRepository{db: db}
}

func (r *GormRepository) List(ctx context.Context, limit int) ([]entity.Todo, error) {
	if limit <= 0 {
		limit = 50
	}

	var todos []entity.Todo
	err := r.db.WithContext(ctx).
		Order("created_at desc").
		Limit(limit).
		Find(&todos).Error
	return todos, err
}

func (r *GormRepository) Create(ctx context.Context, todo *entity.Todo) error {
	return r.db.WithContext(ctx).Create(todo).Error
}

func (r *GormRepository) FindByID(ctx context.Context, id uint) (*entity.Todo, error) {
	var todo entity.Todo
	err := r.db.WithContext(ctx).First(&todo, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return &todo, nil
}

func (r *GormRepository) Update(ctx context.Context, todo *entity.Todo) error {
	return r.db.WithContext(ctx).Save(todo).Error
}

func (r *GormRepository) Delete(ctx context.Context, id uint) error {
	result := r.db.WithContext(ctx).Delete(&entity.Todo{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrNotFound
	}
	return nil
}
