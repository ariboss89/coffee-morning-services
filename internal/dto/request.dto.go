package dto

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email" example:"sosmed89@gmail.com"`
	Password string `json:"password" binding:"required,min=8" example:"12345678"`
}

type RegisterRequest struct {
	Email    string `json:"email" binding:"required,email" example:"sosmed89@gmail.com"`
	Password string `json:"password" binding:"required,min=8" example:"sosmed123"`
}
