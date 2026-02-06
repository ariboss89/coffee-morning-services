package service

import (
	"context"
	"errors"
	"log"

	"github.com/ariboss89/coffee-morning-services/internal/dto"
	"github.com/ariboss89/coffee-morning-services/internal/repository"
	"github.com/ariboss89/coffee-morning-services/pkg/hash"
	"github.com/ariboss89/coffee-morning-services/pkg/jwt"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

type AuthService struct {
	authRepository *repository.AuthRepository
	redis          *redis.Client
	db             *pgxpool.Pool
}

func NewAuthService(authRepository *repository.AuthRepository, rdb *redis.Client, db *pgxpool.Pool) *AuthService {
	return &AuthService{authRepository: authRepository, redis: rdb, db: db}
}

// func (a AuthService) Register(ctx context.Context, newUser dto.RegisterRequest) (dto.RegisterResponse, error) {
// 	hc := hash.HashConfig{}
// 	hc.UseRecommended()

// 	hp, err := hc.GenHash(newUser.Password)
// 	if err != nil {
// 		return dto.RegisterResponse{}, err
// 	}
// 	newUser.Password = hp
// 	if err := a.authRepository.Register(ctx, a.db, newUser); err != nil {
// 		return dto.RegisterResponse{}, err
// 	}

// 	response := dto.RegisterResponse{
// 		ResponseSuccess: dto.ResponseSuccess{
// 			Message: "New Account Registered !",
// 			Status:  "Success",
// 		},
// 	}

// 	return response, nil
// }

func (a AuthService) Register(ctx context.Context, newUser dto.RegisterRequest) (dto.RegisterResponse, error) {
	tx, err := a.db.Begin(ctx)

	if err != nil {
		log.Println(err)
		return dto.RegisterResponse{
			ResponseSuccess: dto.ResponseSuccess{},
		}, err
	}

	hc := hash.HashConfig{}
	hc.UseRecommended()

	hp, err := hc.GenHash(newUser.Password)
	if err != nil {
		return dto.RegisterResponse{}, err
	}
	newUser.Password = hp
	data, err := a.authRepository.Register(ctx, tx, newUser)

	if err != nil {
		return dto.RegisterResponse{}, err
	}
	defer tx.Rollback(ctx)

	if err = a.authRepository.InsertUsers(ctx, tx, data); err != nil {
		return dto.RegisterResponse{}, err
	}

	if e := tx.Commit(ctx); e != nil {
		log.Println("failed to commit", e.Error())
		return dto.RegisterResponse{}, e
	}

	response := dto.RegisterResponse{
		ResponseSuccess: dto.ResponseSuccess{
			Message: "New Account Registered !",
			Status:  "Success",
		},
	}

	return response, nil
}

func (a AuthService) Login(ctx context.Context, login dto.LoginRequest) (dto.LoginResponse, error) {
	tx, err := a.db.Begin(ctx)
	if err != nil {
		log.Println(err)
		return dto.LoginResponse{
			ResponseSuccess: dto.ResponseSuccess{},
		}, err
	}

	var resp dto.LoginResponse
	hc := hash.HashConfig{}

	data, err := a.authRepository.Login(ctx, tx, login)
	if err != nil {
		return dto.LoginResponse{}, err
	}

	defer tx.Rollback(ctx)

	hp, err := hc.ComparePwdAndHash(login.Password, data.Password)
	if err != nil {
		return dto.LoginResponse{}, err
	}

	if hp {
		claims := jwt.NewJWTClaims(data.Id, data.Email, data.Role)
		token, err := claims.GenToken()

		if err != nil {
			return dto.LoginResponse{}, err
		}

		if err = a.authRepository.UpdateLastLogin(ctx, tx, data.Id); err != nil {
			return dto.LoginResponse{}, err
		}

		if e := tx.Commit(ctx); e != nil {
			log.Println("failed to commit", e.Error())
			return dto.LoginResponse{}, e
		}

		resp = dto.LoginResponse{
			ResponseSuccess: dto.ResponseSuccess{
				Status:  "Success",
				Message: "Login Success",
			},
			Data: dto.JWT{
				Token: token,
			},
		}

	} else {

		return dto.LoginResponse{}, errors.New("username or password is wrong")
	}

	return resp, nil
}

func (a AuthService) LogoutUser(ctx context.Context, token string) (bool, error) {
	rkey := "ari:tickitz:logout" + token
	//rsc := a.redis.Get(ctx, rkey)

	rsc, err := a.redis.Exists(ctx, rkey).Result()

	if err != nil {
		return false, err
	}

	if rsc == 0 {
		a.redis.Set(ctx, rkey, token, 0)
	}

	return rsc > 0, nil
}
