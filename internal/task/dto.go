package task

import "github.com/google/uuid"

type CreateTaskDTO struct {
	CreatorID uuid.UUID `json:"creator_id" validate:"required"`
	Body      string    `json:"body,omitempty" validate:"max=2048"`
	Title     string    `json:"title" validate:"required,min=1,max=255"`
	ListID    *string   `json:"list_id,omitempty" validate:"uuid4"`
}

type UpdateTaskDTO struct {
	TaskID uuid.UUID `json:"task_id" validate:"required"`
	Title  *string   `json:"title" validate:"min=1,max=255"`
	Body   *string   `json:"body" validate:"min=1,max=2048"`
	Status *bool     `json:"status"`
}

type UpdateFieldsDTO struct {
	NewBody   *string
	NewTitle  *string
	NewStatus *bool
}
