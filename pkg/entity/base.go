package entity

import (
	"time"

	"github.com/google/uuid"
	"gopkg.in/guregu/null.v4"
)

type Base struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	DeletedAt null.Time `json:"deleted_at"`
}

func (b *Base) Before() {
	if b.ID == uuid.Nil {
		b.ID = uuid.New()
		b.CreatedAt = time.Now()
	}
	b.UpdatedAt = time.Now()
}
