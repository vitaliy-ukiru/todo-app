package postgres

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgtype"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/pkg/errors"
	"github.com/vitaliy-ukiru/todo-app/internal/task"
	"github.com/vitaliy-ukiru/todo-app/internal/task/storage/postgres/gen"
	"github.com/vitaliy-ukiru/todo-app/pkg/pgxuuid"
)

type Repository struct {
	q *gen.DBQuerier
	p *pgxpool.Pool
}

func (r Repository) Create(ctx context.Context, newTask *task.Task) error {
	row, err := r.q.CreateTask(ctx, gen.CreateTaskParams{
		CreatorID: pgxuuid.New(newTask.CreatorID),
		ListID:    pgxuuid.NewPointer(newTask.ListID),
		Title:     newTask.Title,
		Body:      newTask.Body,
	})
	if err != nil {
		return errors.WithStack(err)
	}
	newTask.ID = row.ID.UUID
	newTask.CreatedAt = row.CreatedAt.Time
	return nil
}

func (r Repository) ByID(ctx context.Context, id uuid.UUID) (*task.Task, error) {
	row, err := r.q.FindTaskByID(ctx, pgxuuid.New(id))
	if err != nil {
		return nil, errors.WithStack(err)
	}
	t := &task.Task{
		BasicTask: rowToBaseTask(rowType{
			ID:        row.ID,
			Title:     row.Title,
			Body:      row.Body,
			Done:      row.Done,
			CreatedAt: row.CreatedAt,
			UpdatedAt: row.UpdatedAt,
		}),
		CreatorID: row.CreatorID.UUID,
	}
	if row.ListID.Status == pgtype.Present {
		t.ListID = &row.ListID.UUID
	}

	return t, nil

}

func (r Repository) InOneListBasic(ctx context.Context, list uuid.UUID) ([]task.BasicTask, error) {
	rows, err := r.q.FindBasicTaskInList(ctx, pgxuuid.New(list))
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return convertSlice(rows, func(in gen.FindBasicTaskInListRow) task.BasicTask {
		return rowToBaseTask(rowType(in))
	}), nil
}

func (r Repository) InNullList(ctx context.Context, user uuid.UUID) ([]task.BasicTask, error) {
	list, err := r.q.FindTaskInMainList(ctx, pgxuuid.New(user))
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return convertSlice(list, func(in gen.FindTaskInMainListRow) task.BasicTask {
		return rowToBaseTask(rowType(in))
	}), nil

}

func (r Repository) Update(ctx context.Context, id uuid.UUID, title, body *string, status *bool) error {
	return errors.WithStack(r.p.BeginFunc(ctx, func(tx pgx.Tx) error {
		q := gen.NewQuerier(tx)
		pgId := pgxuuid.New(id)
		if title != nil {
			if _, err := q.UpdateTaskTitle(ctx, *title, pgId); err != nil {
				return errors.WithStack(err)
			}
		}

		if body != nil {
			if _, err := q.UpdateTaskBody(ctx, *body, pgId); err != nil {
				return errors.WithStack(err)
			}
		}

		if status != nil {
			if _, err := q.UpdateTaskStatus(ctx, *status, pgId); err != nil {
				return errors.WithStack(err)
			}
		}

		return nil
	}))
}

func (r Repository) UpdateStatus(ctx context.Context, taskId uuid.UUID) (bool, error) {
	status, err := r.q.ChangeTaskStatus(ctx, pgxuuid.New(taskId))
	return status, errors.WithStack(err)
}

func (r Repository) Delete(ctx context.Context, taskId uuid.UUID) error {
	_, err := r.q.DeleteTask(ctx, pgxuuid.New(taskId))
	return errors.WithStack(err)
}

func (r Repository) Ping(ctx context.Context) error {
	return r.p.Ping(ctx)
}

type rowType struct {
	ID        pgxuuid.UUID       `json:"id"`
	Title     string             `json:"title"`
	Body      string             `json:"body"`
	Done      bool               `json:"done"`
	CreatedAt pgtype.Timestamptz `json:"created_at"`
	UpdatedAt pgtype.Timestamptz `json:"updated_at"`
}

func rowToBaseTask(row rowType) task.BasicTask {
	t := task.BasicTask{
		Body:      row.Body,
		CreatedAt: row.CreatedAt.Time,
		ID:        row.ID.UUID,
		Title:     row.Title,
		Done:      row.Done,
	}

	if row.UpdatedAt.Status == pgtype.Present {
		t.UpdateAt = &row.UpdatedAt.Time
	}
	return t
}

func convertSlice[In, Out any](input []In, convertFn func(In) Out) []Out {
	result := make([]Out, len(input))
	for i, item := range input {
		result[i] = convertFn(item)
	}
	return result
}
