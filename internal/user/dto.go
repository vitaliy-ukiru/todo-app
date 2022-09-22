package user

import "github.com/google/uuid"

type UpdatePasswordUserDTO struct {
	UserID         uuid.UUID `json:"user_id" validate:"required"`
	ActualPassword string    `json:"password" validate:"required,min=6"`
	NewPassword    string    `json:"new_password" validate:"required,min=6"`
}

type CreateUserDTO struct {
	Email    string `json:"email" validate:"required,email"`
	Username string `json:"username" validate:"required,min=1,max=255"`
	Password string `json:"password" validate:"required, min=6"`
}
