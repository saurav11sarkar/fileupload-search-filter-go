package user

type UpdateUserRequest struct {
	Name     *string `json:"name" validate:"omitempty,min=3,max=100"`
	Email    *string `json:"email" validate:"omitempty,email"`
	Password *string `json:"password" validate:"omitempty,min=6,max=255"`
	Role     *string `json:"role" validate:"omitempty"`
	Status   *string `json:"status" validate:"omitempty"`
}

type UpdateProfileRequest struct {
	Profile string `json:"profile" validate:"required"`
}

type UserResponse struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	Password  string `json:"password"`
	Role      string `json:"role"`
	Status    string `json:"status"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}
