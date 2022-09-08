package task

type CreateTaskDTO struct {
	CreatorID string  `json:"creator_id" validate:"required,uui4"`
	Body      string  `json:"body,omitempty" validate:"max=2048"`
	Title     string  `json:"title" validate:"required,min=1,max=255"`
	ListID    *string `json:"list_id,omitempty" validate:"uuid4"`
}

type UpdateTaskDTO struct {
	TaskID string  `json:"task_id" validate:"required,uuid4"`
	Title  *string `json:"title" validate:"min=1,max=255"`
	Body   *string `json:"body" validate:"min=1,max=2048"`
	Status *bool   `json:"status"`
}

type UpdateFieldsDTO struct {
	NewBody   *string
	NewTitle  *string
	NewStatus *bool
}
