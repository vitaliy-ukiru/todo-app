package postgres

import (
	"context"
	"database/sql"
	"database/sql/driver"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	domain "github.com/vitaliy-ukiru/todo-app/internal/user"
	"github.com/vitaliy-ukiru/todo-app/pkg/pgxuuid"
)

type Repository struct {
	q *DBQuerier
	c Connection
}

type Connection interface {
	genericConn
	driver.Pinger
}

func NewRepository(c Connection) *Repository {
	return &Repository{c: c, q: NewQuerier(c)}
}

func (r Repository) Create(ctx context.Context, user *domain.User) error {
	userRow, err := r.q.CreateUser(ctx, CreateUserParams{
		Email:    user.Email,
		Username: user.Username,
		Password: user.Password,
	})
	if err != nil {
		return errors.WithStack(err)
	}
	user.ID = userRow.ID.UUID
	user.CreatedAt = userRow.CreatedAt.Time
	return nil
}

func (r Repository) FindByID(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	row, err := r.q.FindUserByID(ctx, pgxuuid.New(id))
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return &domain.User{
		ID:        row.ID.UUID,
		CreatedAt: row.CreatedAt.Time,
		Email:     row.Email,
		Username:  row.Username,
	}, nil
}

func (r Repository) FindByEmail(ctx context.Context, email string) (*domain.User, error) {
	row, err := r.q.FindUserByEmail(ctx, email)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return &domain.User{
		ID:        row.ID.UUID,
		CreatedAt: row.CreatedAt.Time,
		Email:     row.Email,
		Username:  row.Username,
		Password:  row.Password,
	}, nil
}

func (r Repository) Exists(ctx context.Context, id uuid.UUID) error {
	count, err := r.q.ExistsUser(ctx, pgxuuid.New(id))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.ErrUserNotFound
		}

		return errors.WithStack(err)
	}
	if count == nil || *count == 0 {
		return domain.ErrUserNotFound
	}

	return nil
}

func (r Repository) DeleteUser(ctx context.Context, userID uuid.UUID) error {
	_, err := r.q.DeleteUser(ctx, pgxuuid.New(userID))
	return errors.WithStack(err)
}

func (r Repository) UpdatePassword(ctx context.Context, userID uuid.UUID, newPassword string) error {
	_, err := r.q.UpdateUserPassword(ctx, newPassword, pgxuuid.New(userID))
	return errors.WithStack(err)
}

func (r Repository) Ping(ctx context.Context) error {
	return r.c.Ping(ctx)
}
