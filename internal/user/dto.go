package user

type UpdatePasswordUserDTO struct {
	UserID         string `json:"user_id"`
	ActualPassword string `json:"password,omitempty"`
	NewPassword    string `json:"new_password,omitempty"`
}

type CreateUserDTO struct {
	Email    string `json:"email" validate:"required,email"`
	Username string `json:"username" validate:"required,min=1,max=255"`
	Password string `json:"password" validate:"required, min=6"`
}
