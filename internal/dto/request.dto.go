package dto

import "mime/multipart"

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email" example:"sosmed89@gmail.com"`
	Password string `json:"password" binding:"required,min=8" example:"12345678"`
}

type RegisterRequest struct {
	Email    string `json:"email" binding:"required,email" example:"sosmed89@gmail.com"`
	Password string `json:"password" binding:"required,min=8" example:"sosmed123"`
}

type UserRequest struct {
	Fullname   string                `form:"fullname,omitempty" json:"fullname"`
	AvatarFile *multipart.FileHeader `form:"avatar_file,omitempty" json:"avatar_file"`
	Avatar     string                `form:"avatar,omitempty" json:"avatar"`
	Bio        string                `form:"bio,omitempty" json:"bio"`
}

type InteractionRequest struct {
	ContentFile *multipart.FileHeader `form:"content_file,omitempty" json:"content_file"`
	ContentName string                `form:"content_name,omitempty" json:"content_name"`
	Caption     string                `form:"caption,omitempty" json:"caption"`
}
