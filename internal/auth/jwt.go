package auth

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/vitaliy-ukiru/todo-app/internal/user"
	"github.com/vitaliy-ukiru/todo-app/pkg/jwt"
)

type Claims struct {
	ID       uuid.UUID `json:"id"` // uuid
	Username string    `json:"username"`

	jwt.RegisteredClaims
}

type UserDTO struct {
	ID        uuid.UUID `json:"id,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty"`
	Email     string    `json:"email"`
	Username  string    `json:"username"`
}

type UserChecker interface {
	CheckCredentials(ctx context.Context, email string, password string) (*user.User, error)
}

type JwtService struct {
	provider        *jwt.Provider
	accessTokenLife time.Duration
	userUC          UserChecker
}

func NewJwtService(provider *jwt.Provider, accessTokenLife time.Duration) *JwtService {
	return &JwtService{provider: provider, accessTokenLife: accessTokenLife}
}

type CredentialsDTO struct {
	Email    string
	Password string
}

type LoginResult struct {
	Token     string    `json:"access_token"`
	ExpiresAt time.Time `json:"expires_at"`
}

func (j JwtService) Login(ctx context.Context, cred CredentialsDTO) (LoginResult, error) {
	u, err := j.userUC.CheckCredentials(ctx, cred.Email, cred.Password)
	if err != nil {
		return LoginResult{}, errors.WithStack(err)
	}
	token, expiresAt, err := j.NewAccessToken(*u)
	if err != nil {
		return LoginResult{}, errors.WithStack(err)
	}

	return LoginResult{Token: token, ExpiresAt: expiresAt}, nil
}

func (j JwtService) NewAccessToken(u user.User) (string, time.Time, error) {
	claims := &Claims{
		ID:       u.ID,
		Username: u.Username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.AddDurationToNow(j.accessTokenLife),
		},
	}
	token, err := j.provider.CreateToken(claims)
	return token, claims.ExpiresAt.Time, errors.WithStack(err)
}

func (j JwtService) VerifyToken(token string) (*Claims, error) {
	claims := new(Claims)
	err := j.provider.VerifyToken(token, claims)
	return claims, errors.WithStack(err)
}
