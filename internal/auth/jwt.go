package auth

import (
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

type JwtService struct {
	provider        *jwt.Provider
	accessTokenLife time.Duration
}

func NewJwtService(provider *jwt.Provider, accessTokenLife time.Duration) *JwtService {
	return &JwtService{provider: provider, accessTokenLife: accessTokenLife}
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
