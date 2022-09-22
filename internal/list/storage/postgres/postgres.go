package postgres

import (
	"context"
	"database/sql/driver"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/pkg/errors"
	"github.com/vitaliy-ukiru/todo-app/internal/list"
	"github.com/vitaliy-ukiru/todo-app/pkg/pgxuuid"
)

type Repository struct {
	q *DBQuerier
	p *pgxpool.Pool
}

func NewRepository(c Connection) *Repository {
	return &Repository{c: c, q: NewQuerier(c)}
}

func (r Repository) CreateList(ctx context.Context, title string, creator uuid.UUID) (uuid.UUID, error) {
	listId, err := r.q.CreateList(ctx, pgxuuid.New(creator), title)
	return listId.UUID, errors.WithStack(err)

}

func (r Repository) FindByID(ctx context.Context, listId uuid.UUID) (*list.TaskListInfo, error) {
	listRow, err := r.q.FindListByID(ctx, pgxuuid.New(listId))
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return &list.TaskListInfo{
		CreatorID: listRow.CreatorID.UUID,
		ID:        listRow.ID.UUID,
		Title:     listRow.Title,
	}, nil
}

func (r Repository) FindByUserID(ctx context.Context, userId uuid.UUID) ([]list.TaskListInfo, error) {
	lists, err := r.q.FindUserLists(ctx, pgxuuid.New(userId))
	if err != nil {
		return nil, errors.WithStack(err)
	}
	result := make([]list.TaskListInfo, len(lists))
	for i, listRow := range lists {
		result[i] = list.TaskListInfo{
			CreatorID: listRow.CreatorID.UUID,
			ID:        listRow.ID.UUID,
			Title:     listRow.Title,
		}
	}
	return result, nil

}

func (r Repository) UpdateTitle(ctx context.Context, listId uuid.UUID, newTitle string) error {
	_, err := r.q.UpdateListTitle(ctx, newTitle, pgxuuid.New(listId))
	//TODO: add check on rows affected
	return errors.WithStack(err)
}

func (r Repository) Delete(ctx context.Context, listId uuid.UUID) error {
	_, err := r.q.DeleteList(ctx, pgxuuid.New(listId))
	return errors.WithStack(err)
}

func (r Repository) Ping(ctx context.Context) error {
	return r.p.Ping(ctx)
}
