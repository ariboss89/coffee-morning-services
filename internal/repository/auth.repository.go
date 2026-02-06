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
	Login(ctx context.Context, req dto.LoginRequest, db DBTX) (model.User, error)
	UpdateLastLogin(ctx context.Context, db DBTX, id int) error
}

type AuthRepository struct{}

func NewAuthRepository() *AuthRepository {
	return &AuthRepository{}
}

func (ar *AuthRepository) Login(ctx context.Context, db DBTX, req dto.LoginRequest) (model.User, error) {
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

	var user model.User
	err := row.Scan(
		&user.Id,
		&user.Email,
		&user.Password,
		&user.Role,
	)

	if err != nil {
		log.Println(err.Error())
		if errors.Is(err, pgx.ErrNoRows) {
			return model.User{}, err
		}
		return model.User{}, err
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

func (ar *AuthRepository) Register(ctx context.Context, db DBTX, req dto.RegisterRequest) error {
	query := `
		INSERT INTO
		    accounts (email, password)
		VALUES
		    ($1, $2)
	`

	_, err := db.Exec(ctx, query, req.Email, req.Password)
	if err != nil {
		log.Println(err.Error())
		if strings.Contains(err.Error(), "duplicate") {
			return err
		}
		return err
	}

	return nil
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
