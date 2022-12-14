package task

import (
	"time"

	"github.com/google/uuid"
)

type BasicTask struct {
	ID        uuid.UUID  `json:"id,omitempty"`
	Title     string     `json:"title"`
	Body      string     `json:"body,omitempty"`
	CreatedAt time.Time  `json:"created_at,omitempty"`
	UpdateAt  *time.Time `json:"update_at,omitempty"`
	Done      bool       `json:"done"`
}

type Task struct {
	CreatorID uuid.UUID `json:"creator_id"`

	ListID *uuid.UUID `json:"list_id,omitempty"`

	BasicTask
}
