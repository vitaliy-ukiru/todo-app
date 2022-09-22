package user

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID `json:"id,omitempty" db:"id"`
	CreatedAt time.Time `json:"created_at,omitempty" db:"created_at"`
	Email     string    `json:"email" db:"email"`
	Username  string    `json:"username" db:"username"`

	Password string `json:"-" db:"password"`
}
