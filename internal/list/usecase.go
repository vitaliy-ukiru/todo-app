package list

import (
	"context"
	"html"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/pkg/errors"
)

type Usecase interface {
	Create(ctx context.Context, list CreateTaskListDTO) (*TaskListInfo, error)

	FindById(ctx context.Context, listId uuid.UUID) (*TaskListInfo, error)
	FindUserLists(ctx context.Context, userId uuid.UUID) ([]TaskListInfo, error)
	FullTaskList(ctx context.Context, listId uuid.UUID) (*TaskList, error)

	UpdateTitle(ctx context.Context, list UpdateTaskListDTO) (*TaskListInfo, error)

	Delete(ctx context.Context, listId uuid.UUID) error
	Ping(ctx context.Context) (time.Duration, error)
}

type Storage interface {
	CreateList(ctx context.Context, title string, creator uuid.UUID) (uuid.UUID, error)
	FindByID(ctx context.Context, listId uuid.UUID) (*TaskListInfo, error)
	FindByUserID(ctx context.Context, userId uuid.UUID) ([]TaskListInfo, error)
	UpdateTitle(ctx context.Context, listId uuid.UUID, newTitle string) error
	Delete(ctx context.Context, listId uuid.UUID) error

	Ping(ctx context.Context) error
}

type UserUsecase interface {
	Exists(ctx context.Context, userId uuid.UUID) error
}

type TaskUsecase interface {
	FindInList(ctx context.Context, listId uuid.UUID) ([]TaskInfoDTO, error)
}

type Service struct {
	store  Storage
	userUC UserUsecase
	taskUC TaskUsecase
}

func NewService(store Storage, userUC UserUsecase, taskUC TaskUsecase) *Service {
	return &Service{store: store, userUC: userUC, taskUC: taskUC}
}

func (s Service) Create(ctx context.Context, dto CreateTaskListDTO) (*TaskListInfo, error) {
	ctxUser, cancelUser := context.WithTimeout(ctx, 2*time.Second)
	defer cancelUser()

	if err := s.userUC.Exists(ctxUser, dto.CreatorID); err != nil {
		return nil, errors.WithStack(err)
	}

	dto.Title = html.EscapeString(strings.TrimSpace(dto.Title))
	listId, err := s.store.CreateList(ctx, dto.Title, dto.CreatorID)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return &TaskListInfo{
		CreatorID: dto.CreatorID,
		ID:        listId,
		Title:     dto.Title,
	}, nil
}

func (s Service) FindById(ctx context.Context, listId uuid.UUID) (*TaskListInfo, error) {
	listInfo, err := s.store.FindByID(ctx, listId)
	return listInfo, errors.WithStack(err)
}

func (s Service) FindUserLists(ctx context.Context, userId uuid.UUID) ([]TaskListInfo, error) {
	listsInfo, err := s.store.FindByUserID(ctx, userId)
	return listsInfo, errors.WithStack(err)
}

func (s Service) FullTaskList(ctx context.Context, listId uuid.UUID) (*TaskList, error) {
	listInfo, err := s.FindById(ctx, listId)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	tasks, err := s.taskUC.FindInList(ctx, listId)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return &TaskList{
		TaskListInfo: *listInfo,
		Tasks:        tasks,
	}, nil
}

func (s Service) UpdateTitle(ctx context.Context, list UpdateTaskListDTO) (*TaskListInfo, error) {
	title := html.EscapeString(strings.TrimSpace(list.NewTitle))

	if err := s.store.UpdateTitle(ctx, list.ListID, title); err != nil {
		return nil, errors.WithStack(err)
	}

	taskListInfo, err := s.FindById(ctx, list.ListID)
	return taskListInfo, errors.WithStack(err)
}

func (s Service) Delete(ctx context.Context, listId uuid.UUID) error {
	return errors.WithStack(s.store.Delete(ctx, listId))
}

func (s Service) Ping(ctx context.Context) (time.Duration, error) {
	start := time.Now()
	err := s.store.Ping(ctx)
	ping := time.Since(start)

	return ping, errors.WithStack(err)
}
