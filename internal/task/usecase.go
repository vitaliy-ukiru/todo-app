package task

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/pkg/errors"
)

type Usecase interface {
	Create(ctx context.Context, task CreateTaskDTO) (*Task, error)

	FindByID(ctx context.Context, taskId uuid.UUID) (*Task, error)
	FindInList(ctx context.Context, listId uuid.UUID) ([]BasicTask, error)
	MainList(ctx context.Context, userId uuid.UUID) ([]BasicTask, error)

	ChangeStatus(ctx context.Context, taskId uuid.UUID) (bool, error)
	UpdateTask(ctx context.Context, task UpdateTaskDTO) (*Task, error)

	Delete(ctx context.Context, taskId uuid.UUID) error
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
	Exists(ctx context.Context, userId uuid.UUID) error
}

type Service struct {
	store  Storage
	userUC UserUsecase
}

func (s Service) Create(ctx context.Context, dto CreateTaskDTO) (*Task, error) {
	if err := s.userUC.Exists(ctx, dto.CreatorID); err != nil {
		return nil, errors.WithStack(err)
	}

	newTask := &Task{
		CreatorID: dto.CreatorID,
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

func (s Service) FindByID(ctx context.Context, taskId uuid.UUID) (*Task, error) {
	task, err := s.store.ByID(ctx, taskId)
	return task, errors.WithStack(err)
}

func (s Service) FindInList(ctx context.Context, listId uuid.UUID) ([]BasicTask, error) {
	tasks, err := s.store.InOneListBasic(ctx, listId)
	return tasks, errors.WithStack(err)
}

func (s Service) MainList(ctx context.Context, userId uuid.UUID) ([]BasicTask, error) {
	tasks, err := s.store.InNullList(ctx, userId)
	return tasks, errors.WithStack(err)
}

func (s Service) ChangeStatus(ctx context.Context, taskId uuid.UUID) (bool, error) {
	status, err := s.store.UpdateStatus(ctx, taskId)
	return status, errors.WithStack(err)

}

func (s Service) UpdateTask(ctx context.Context, update UpdateTaskDTO) (*Task, error) {
	if err := s.store.Update(ctx, update.TaskID, update.Title, update.Body, update.Status); err != nil {
		return nil, errors.WithStack(err)
	}

	task, err := s.store.ByID(ctx, update.TaskID)
	return task, errors.WithStack(err)
}

func (s Service) Delete(ctx context.Context, taskId uuid.UUID) error {
	return errors.WithStack(s.store.Delete(ctx, taskId))
}

func (s Service) Ping(ctx context.Context) (time.Duration, error) {
	start := time.Now()
	err := s.store.Ping(ctx)
	ping := time.Since(start)

	return ping, errors.WithStack(err)
}
