package list

import (
	"time"

	"github.com/google/uuid"
)

type CreateTaskListDTO struct {
	CreatorID uuid.UUID `json:"creator_id" validate:"required"`
	Title     string    `json:"title" validate:"required,min=1,max=255"`
}

type UpdateTaskListDTO struct {
	ListID   uuid.UUID `json:"list_id" validate:"required"`
	NewTitle string    `json:"new_title" validate:"required,min=1,max=255"`
}

type TaskInfoDTO struct {
	Body      string     `json:"body"`
	CreatedAt time.Time  `json:"created_at"`
	ID        uuid.UUID  `json:"id"`
	Title     string     `json:"title"`
	UpdateAt  *time.Time `json:"update_at,omitempty"`
	Done      bool       `json:"done"`
}
