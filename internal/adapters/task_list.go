package adapters

import (
	"context"

	"github.com/google/uuid"
	"github.com/vitaliy-ukiru/todo-app/internal/list"
	"github.com/vitaliy-ukiru/todo-app/internal/task"
)

type TaskUsecaseAdaptor struct {
	task.Usecase
}

func NewTaskUCAdaptor(usecase task.Usecase) list.TaskUsecase {
	return &TaskUsecaseAdaptor{Usecase: usecase}
}

func (t TaskUsecaseAdaptor) FindInList(ctx context.Context, listId uuid.UUID) ([]list.TaskInfoDTO, error) {
	tasks, err := t.Usecase.FindInList(ctx, listId)
	if err != nil {
		return nil, err
	}
	result := make([]list.TaskInfoDTO, len(tasks))
	for i, basicTask := range tasks {
		result[i] = list.TaskInfoDTO(basicTask)
	}
	return result, nil
}
