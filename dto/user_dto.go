package dto

type CreateUserRequest struct {
	Username string `json:"username" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
	FullName string `json:"full_name" binding:"required"`
	Phone    string `json:"phone" binding:"omitempty"`
	Position string `json:"position" binding:"omitempty"`
}

type LoginRequest struct {
	Username string `json:"username" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

type LogoutRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

type UpdateUserRequest struct {
	FullName string `json:"full_name" binding:"required"`
	Phone    string `json:"phone" binding:"omitempty"`
	Position string `json:"position" binding:"omitempty"`
}

type UserResponse struct {
	ID        uint   `json:"id"`
	Uuid      string `json:"uuid"`
	Username  string `json:"username"`
	FullName  string `json:"fullname"`
	Phone     string `json:"phone"`
	Position  string `json:"position"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}
