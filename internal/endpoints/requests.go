package endpoints

import (
	"github.com/aintsashqa/simple-todo/pkg/filter"
	"github.com/google/uuid"
	"gopkg.in/guregu/null.v4"
)

type CreateRequest struct {
	Title       string      `json:"title"`
	Description null.String `json:"description"`
}

type UpdateRequest struct {
	ID          uuid.UUID   `json:"-"`
	Title       string      `json:"title"`
	Description null.String `json:"description"`
}

type GetListRequest struct {
	FilterType  filter.Type `json:"-"`
	CompletedAt null.Time   `json:"-"`
}

type ChangeToCompletedRequest struct {
	ID uuid.UUID `json:"-"`
}
