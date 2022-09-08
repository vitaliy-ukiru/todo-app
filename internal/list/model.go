package list

import (
	"github.com/google/uuid"
)

type TaskListInfo struct {

	// id
	// Format: uuid
	ID uuid.UUID `json:"id,omitempty"`

	// creator id
	// Format: uuid
	CreatorID uuid.UUID `json:"creator_id,omitempty"`

	// title
	// Required: true
	Title string `json:"title"`
}

type TaskList struct {
	TaskListInfo

	// tasks
	Tasks []TaskInfoDTO `json:"tasks"`
}
