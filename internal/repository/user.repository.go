package repository

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/ariboss89/coffee-morning-services/internal/dto"
	"github.com/ariboss89/coffee-morning-services/internal/model"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository struct {
	db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

func (u UserRepository) GetUserProfileById(ctx context.Context, idUser int) (model.User, error) {
	sqlStr := "SELECT fullname, avatar, bio FROM users WHERE user_id = $1"
	rows := u.db.QueryRow(ctx, sqlStr, idUser)
	var user model.User
	if err := rows.Scan(&user.Fullname, &user.Avatar, &user.Bio); err != nil {
		log.Println(err.Error())
		return model.User{}, err
	}
	return user, nil
}

func (u UserRepository) UpdateProfile(ctx context.Context, update dto.UserRequest, userId int) (pgconn.CommandTag, error) {
	var sql strings.Builder
	values := []any{}

	sql.WriteString("UPDATE users SET")
	if update.Fullname != "" {
		fmt.Fprintf(&sql, " fullname=$%d", len(values)+1)
		values = append(values, update.Fullname)
	}
	if update.Bio != "" {
		if len(values) > 0 {
			sql.WriteString(",")
		}
		fmt.Fprintf(&sql, " bio=$%d", len(values)+1)
		values = append(values, update.Bio)
	}
	if update.Avatar != "" {
		if len(values) > 0 {
			sql.WriteString(",")
		}
		fmt.Fprintf(&sql, " avatar=$%d", len(values)+1)
		values = append(values, fmt.Sprintf("/profile/%s", update.Avatar))
	}
	if update.Fullname != "" || update.Bio != "" || update.Avatar != "" {
		sql.WriteString(", updated_at= NOW() WHERE ")
		fmt.Fprintf(&sql, "user_id='%d'", userId)
	}

	sqlStr := sql.String()

	return u.db.Exec(ctx, sqlStr, values...)
}
