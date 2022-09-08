package list

import (
	"time"

	"github.com/google/uuid"
)

type CreateTaskListDTO struct {
	CreatorID string `json:"creator_id" validate:"required,uuid4"`
	Title     string `json:"title" validate:"required,min=1,max=255"`
}

type UpdateTaskListDTO struct {
	ListID   string `json:"list_id" validate:"required,uuid4"`
	NewTitle string `json:"new_title" validate:"required,min=1,max=255"`
}

type TaskInfoDTO struct {
	Body      string     `json:"body"`
	CreatedAt time.Time  `json:"created_at"`
	ID        uuid.UUID  `json:"id"`
	Title     string     `json:"title"`
	UpdateAt  *time.Time `json:"update_at,omitempty"`
	Done      bool       `json:"done"`
}
