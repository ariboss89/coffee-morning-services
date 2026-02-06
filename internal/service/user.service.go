package service

import (
	"context"
	"errors"

	"github.com/ariboss89/coffee-morning-services/internal/dto"
	"github.com/ariboss89/coffee-morning-services/internal/repository"
)

type UserRepo interface {
	GetUserProfileByEmail(ctx context.Context, email string)
}

type UserService struct {
	userRepository *repository.UserRepository
}

func NewUserService(userRepository *repository.UserRepository) *UserService {
	return &UserService{
		userRepository: userRepository,
	}
}

func (u UserService) GetUserProfileById(ctx context.Context, idUser int) (dto.UserResponse, error) {
	data, err := u.userRepository.GetUserProfileById(ctx, idUser)
	if err != nil {
		return dto.UserResponse{}, err
	}

	response := dto.UserResponse{
		Fullname: data.Fullname,
		Avatar:   data.Avatar,
		Bio:      data.Bio,
	}

	return response, nil
}

func (u UserService) UpdateProfile(ctx context.Context, update dto.UserRequest, idUser int) error {
	cmd, err := u.userRepository.UpdateProfile(ctx, update, idUser)
	if err != nil {
		return nil
	}
	if cmd.RowsAffected() == 0 {
		return errors.New("no data updated")
	}
	// invalidasi cache
	return nil
}
