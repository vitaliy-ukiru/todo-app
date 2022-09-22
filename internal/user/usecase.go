package user

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	bcryptx "github.com/vitaliy-ukiru/todo-app/pkg/bcrypt"
)

type Usecase interface {
	Create(ctx context.Context, user CreateUserDTO) (*User, error)

	ByID(ctx context.Context, id uuid.UUID) (*User, error)
	ByEmail(ctx context.Context, email string) (*User, error)

	Exists(ctx context.Context, userId uuid.UUID) error
	CheckCredentials(ctx context.Context, email string, password string) (*User, error)

	UpdatePassword(ctx context.Context, user UpdatePasswordUserDTO) error

	Delete(ctx context.Context, userId uuid.UUID) error
	Ping(ctx context.Context) (time.Duration, error)
}

type Storage interface {
	Create(ctx context.Context, user *User) error

	FindByID(ctx context.Context, id uuid.UUID) (*User, error)
	Exists(ctx context.Context, id uuid.UUID) error
	FindByEmail(ctx context.Context, email string) (*User, error)

	DeleteUser(ctx context.Context, userID uuid.UUID) error
	UpdatePassword(ctx context.Context, userID uuid.UUID, newPassword string) error
	Ping(ctx context.Context) error
}

type Service struct {
	store Storage
}

func NewService(store Storage) *Service {
	return &Service{store: store}
}

func (s Service) Create(ctx context.Context, dto CreateUserDTO) (*User, error) {
	password, err := bcryptx.GenerateHash(dto.Password)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	user := &User{
		Email:    dto.Email,
		Username: dto.Username,
		Password: password,
	}

	if err := s.store.Create(ctx, user); err != nil {
		return nil, errors.WithStack(err)
	}

	return user, nil

}

func (s Service) CheckCredentials(ctx context.Context, email string, password string) (*User, error) {
	user, err := s.ByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, ErrUserNotFound) {
			return nil, ErrInvalidCredentials
		}

		return nil, errors.WithStack(err)
	}

	if err := bcryptx.CompareHash(user.Password, password); err != nil {
		return nil, ErrInvalidCredentials
	}

	user.Password = ""
	return user, nil
}

func (s Service) Exists(ctx context.Context, userId uuid.UUID) error {
	return errors.WithStack(s.store.Exists(ctx, userId))
}

func (s Service) ByID(ctx context.Context, userId uuid.UUID) (*User, error) {
	user, err := s.store.FindByID(ctx, userId)
	return user, errors.WithStack(err)
}

func (s Service) ByEmail(ctx context.Context, email string) (*User, error) {
	user, err := s.store.FindByEmail(ctx, email)
	return user, errors.WithStack(err)
}

func (s Service) UpdatePassword(ctx context.Context, dto UpdatePasswordUserDTO) error {
	user, err := s.store.FindByID(ctx, dto.UserID)
	if err != nil {
		return errors.WithStack(err) //TODO err user not found
	}

	if err := bcryptx.CompareHash(user.Password, dto.ActualPassword); err != nil {
		return errors.WithStack(err)
	}

	// user.Password and update.ActualPassword equals here (view compare up)
	if dto.NewPassword == dto.ActualPassword {
		return errors.New("passwords equals") // TODO: move error to errors.go
	}

	newPassword, err := bcryptx.GenerateHash(dto.NewPassword)
	if err != nil {
		return errors.WithStack(err)
	}

	return errors.WithStack(s.store.UpdatePassword(ctx, user.ID, newPassword))
}

func (s Service) Delete(ctx context.Context, userId uuid.UUID) error {
	return errors.WithStack(s.store.DeleteUser(ctx, userId))
}

func (s Service) Ping(ctx context.Context) (time.Duration, error) {
	start := time.Now()
	err := s.store.Ping(ctx)
	ping := time.Since(start)

	return ping, errors.WithStack(err)
}
