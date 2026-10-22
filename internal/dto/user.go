package dto

type CreateUserRequest struct {
	Username string `json:"username" binding:"required" example:"janedoe"`
	Email    string `json:"email" binding:"required" example:"jane@example.com"`
	Password string `json:"password" binding:"required" example:"example-password"`
}

type UserDetailsResponse struct {
	ID       string `json:"id" binding:"required" example:"550e8400-e29b-41d4-a716-446655440000"`
	Username string `json:"username" binding:"required" example:"janedoe"`
	Email    string `json:"email" binding:"required" example:"jane@example.com"`
}

type LoginUserRequest struct {
	Email    string `json:"email" binding:"required" example:"jane@example.com"`
	Password string `json:"password" binding:"required" example:"example-password"`
}

type LoginResponse struct {
	Token string `json:"token" binding:"required" example:"eyJhbGciOi..."`
}
