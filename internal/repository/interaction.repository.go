package repository

import (
	"context"
	"log"

	"github.com/ariboss89/coffee-morning-services/internal/dto"
	"github.com/jackc/pgx/v5/pgxpool"
)

type InteractionRepository struct {
	db *pgxpool.Pool
}

func NewInteractionRepository(db *pgxpool.Pool) *InteractionRepository {
	return &InteractionRepository{
		db: db,
	}
}

func (ir *InteractionRepository) PostContent(ctx context.Context, content dto.InteractionRequest, idUser int) error {
	query := `
		INSERT INTO
		    posts (description, image, user_id)
		VALUES
		    ($1, $2, $3)
	`

	_, err := ir.db.Exec(ctx, query, content.Caption, content.ContentName, idUser)
	if err != nil {
		log.Println(err.Error())
		return err
	}

	return nil
}
