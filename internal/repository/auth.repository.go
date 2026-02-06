package repository

import (
	"context"
	"errors"
	"log"
	"strings"

	"github.com/ariboss89/coffee-morning-services/internal/dto"
	"github.com/ariboss89/coffee-morning-services/internal/model"
	"github.com/jackc/pgx/v5"
)

type AuthRepo interface {
	Login(ctx context.Context, req dto.LoginRequest, db DBTX) (model.Login, error)
	UpdateLastLogin(ctx context.Context, db DBTX, id int) error
}

type AuthRepository struct{}

func NewAuthRepository() *AuthRepository {
	return &AuthRepository{}
}

func (ar *AuthRepository) Login(ctx context.Context, db DBTX, req dto.LoginRequest) (model.Login, error) {
	query := `
		SELECT
		    id,
		    email,
		    password,
		    role
		FROM
		    accounts
		WHERE
		    email = $1;
	`

	row := db.QueryRow(ctx, query, req.Email)

	var user model.Login
	err := row.Scan(
		&user.Id,
		&user.Email,
		&user.Password,
		&user.Role,
	)

	if err != nil {
		log.Println(err.Error())
		if errors.Is(err, pgx.ErrNoRows) {
			return model.Login{}, err
		}
		return model.Login{}, err
	}

	return user, nil
}

func (ar *AuthRepository) UpdateLastLogin(ctx context.Context, db DBTX, id int) error {
	query := `
		UPDATE accounts
		SET
		    lastlogin_at = NOW()
		WHERE
		    id = $1;
	`

	_, err := db.Exec(ctx, query, id)
	if err != nil {
		log.Println(err.Error())
		return err
	}

	return nil
}

func (ar *AuthRepository) Register(ctx context.Context, db DBTX, req dto.RegisterRequest) (int, error) {
	var idUser int

	query := `
		INSERT INTO
		    accounts (email, password)
		VALUES
		    ($1, $2)
		RETURNING id
	`

	err := db.QueryRow(ctx, query, req.Email, req.Password).Scan(&idUser)
	if err != nil {
		log.Println(err.Error())
		if strings.Contains(err.Error(), "duplicate") {
			return 0, err
		}
		return 0, err
	}

	return idUser, nil
}

func (ar *AuthRepository) CheckEmailExists(ctx context.Context, db DBTX, email string) error {
	query := "SELECT email FROM users WHERE email = $1 AND deleted_at IS NULL"

	var emailExist string
	err := db.QueryRow(ctx, query, email).Scan(&emailExist)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return err
		}
		return err
	}

	return nil
}

func (ar *AuthRepository) InsertUsers(ctx context.Context, db DBTX, idUser int) error {
	query := `
		INSERT INTO
		    users (user_id, created_at)
		VALUES
		    ($1, NOW())
	`

	_, err := db.Exec(ctx, query, idUser)
	if err != nil {
		log.Println(err.Error())
		return err
	}

	return nil
}
