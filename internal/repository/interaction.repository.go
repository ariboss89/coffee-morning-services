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

func NewInteractionRepository() *InteractionRepository {
	return &InteractionRepository{}
}

func (ir *InteractionRepository) PostContent(ctx context.Context, content dto.InteractionRequest, idUser int) error {
	query := `
		INSERT INTO
		    posts (description, image, user_id, created_at)
		VALUES
		    ($1, $2, $3, NOW())
	`

	_, err := ir.db.Exec(ctx, query, content.Caption, content.ContentName, idUser)
	if err != nil {
		log.Println(err.Error())
		return err
	}

	return nil
}

func (ir *InteractionRepository) FollowingUser(ctx context.Context, follow dto.FollowingRequest, idUser int) error {
	query := `
		INSERT INTO
		    followers (follower_id, following_id)
		VALUES
		    ($1, $2)
	`

	_, err := ir.db.Exec(ctx, query, follow.Following_Id, idUser)
	if err != nil {
		log.Println(err.Error())
		return err
	}

	return nil
}

func (ir *InteractionRepository) LikePosts(ctx context.Context, db DBTX, like dto.LikeRequest, idUser int) (string, error) {
	query := `
		INSERT INTO
		    likes (post_id, user_id)
		VALUES
		    ($1, $2)
	`

	_, err := db.Exec(ctx, query, like.Post_Id, idUser)
	if err != nil {
		log.Println(err.Error())
		return "", err
	}

	return "Success like posts", nil
}

func (ir *InteractionRepository) DeletePosts(ctx context.Context, db DBTX, like dto.LikeRequest, idUser int) (string, error) {
	query := `
		DELETE FROM likes WHERE post_id = $1
	`

	_, err := db.Exec(ctx, query, like.Post_Id)
	if err != nil {
		log.Println(err.Error())
		return "", err
	}

	return "Success unlike post", nil
}

func (ir *InteractionRepository) ChecksPosts(ctx context.Context, db DBTX, like dto.LikeRequest, idUser int) (int, error) {
	var postId int

	query := `
		SELECT post_id
			FROM likes
			WHERE post_id = $1 AND user_id = $2
	`

	err := db.QueryRow(ctx, query, like.Post_Id, idUser).Scan(&postId)

	if err != nil {
		log.Println(err.Error())
		return 0, err
	}

	return postId, nil
}
