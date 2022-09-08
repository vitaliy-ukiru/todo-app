package task

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/pkg/errors"
)

type Usecase interface {
	Create(ctx context.Context, task CreateTaskDTO) (*Task, error)

	FindByID(ctx context.Context, taskId string) (*Task, error)
	FindInList(ctx context.Context, listId string) ([]BasicTask, error)
	MainList(ctx context.Context, userId string) ([]BasicTask, error)

	ChangeStatus(ctx context.Context, taskId string) (bool, error)
	UpdateTask(ctx context.Context, task UpdateTaskDTO) (*Task, error)

	Delete(ctx context.Context, taskId string) error
	Ping(ctx context.Context) (time.Duration, error)
}

type Storage interface {
	Create(ctx context.Context, task *Task) error

	ByID(ctx context.Context, id uuid.UUID) (*Task, error)
	InOneListBasic(ctx context.Context, list uuid.UUID) ([]BasicTask, error)
	InNullList(ctx context.Context, user uuid.UUID) ([]BasicTask, error)

	Update(ctx context.Context, id uuid.UUID, title, body *string, status *bool) error
	UpdateStatus(ctx context.Context, taskId uuid.UUID) (bool, error)

	Delete(ctx context.Context, taskId uuid.UUID) error

	Ping(ctx context.Context) error
}

type UserUsecase interface {
	Exists(ctx context.Context, userId string) error
}

type Service struct {
	store  Storage
	userUC UserUsecase
}

func (s Service) Create(ctx context.Context, dto CreateTaskDTO) (*Task, error) {
	creatorID, err := uuid.Parse(dto.CreatorID)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	if err := s.userUC.Exists(ctx, dto.CreatorID); err != nil {
		return nil, errors.WithStack(err)
	}

	newTask := &Task{
		CreatorID: creatorID,
		ListID:    nil,
		BasicTask: BasicTask{
			Body:  dto.Body,
			Title: dto.Title,
		},
	}
	if dto.ListID != nil {
		list, err := uuid.Parse(*dto.ListID)
		if err != nil {
			return nil, errors.WithStack(err)
		}
		newTask.ListID = &list
	}

	if err := s.store.Create(ctx, newTask); err != nil {
		return nil, errors.WithStack(err)
	}
	return newTask, nil
}

func (s Service) FindByID(ctx context.Context, id string) (*Task, error) {
	taskId, err := uuid.Parse(id)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	task, err := s.store.ByID(ctx, taskId)
	return task, errors.WithStack(err)
}

func (s Service) FindInList(ctx context.Context, id string) ([]BasicTask, error) {
	listId, err := uuid.Parse(id)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	tasks, err := s.store.InOneListBasic(ctx, listId)
	return tasks, errors.WithStack(err)
}

func (s Service) MainList(ctx context.Context, userID string) ([]BasicTask, error) {
	userId, err := uuid.Parse(userID)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	tasks, err := s.store.InNullList(ctx, userId)
	return tasks, errors.WithStack(err)
}

func (s Service) ChangeStatus(ctx context.Context, taskID string) (bool, error) {
	taskId, err := uuid.Parse(taskID)
	if err != nil {
		return false, errors.WithStack(err)
	}
	status, err := s.store.UpdateStatus(ctx, taskId)
	return status, errors.WithStack(err)

}

func (s Service) UpdateTask(ctx context.Context, update UpdateTaskDTO) (*Task, error) {
	taskId, err := uuid.Parse(update.TaskID)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	if err := s.store.Update(ctx, taskId, update.Title, update.Body, update.Status); err != nil {
		return nil, errors.WithStack(err)
	}

	task, err := s.store.ByID(ctx, taskId)
	return task, errors.WithStack(err)
}

func (s Service) Delete(ctx context.Context, id string) error {
	taskId, err := uuid.Parse(id)
	if err != nil {
		return errors.WithStack(err)
	}

	return errors.WithStack(s.store.Delete(ctx, taskId))
}

func (s Service) Ping(ctx context.Context) (time.Duration, error) {
	start := time.Now()
	err := s.store.Ping(ctx)
	ping := time.Since(start)

	return ping, errors.WithStack(err)
}
