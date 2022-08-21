package internal

import (
	"context"
	"errors"

	"github.com/aintsashqa/simple-todo/pkg/entity"
	"github.com/aintsashqa/simple-todo/pkg/filter"
	"github.com/google/uuid"
	"gopkg.in/guregu/null.v4"
)

var (
	ErrNotFound = errors.New("todo not found")
)

type Todo struct {
	entity.Base
	Title       string      `json:"title"`
	Description null.String `json:"description"`
	CompletedAt null.Time   `json:"completed_at"`
}

type Repository interface {
	Create(ctx context.Context, todo *Todo) error
	Update(ctx context.Context, todo *Todo) error
	GetByID(ctx context.Context, id uuid.UUID) (Todo, error)
	GetList(ctx context.Context, filters ...filter.Filter) ([]Todo, error)
	ChangeCompleteStatus(ctx context.Context, todo *Todo) error
}
