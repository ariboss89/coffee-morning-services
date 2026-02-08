package service

import (
	"context"
	"log"

	"github.com/ariboss89/coffee-morning-services/internal/dto"
	"github.com/ariboss89/coffee-morning-services/internal/repository"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

type InteractionService struct {
	interactionRepository *repository.InteractionRepository
	redis                 *redis.Client
	db                    *pgxpool.Pool
}

func NewInteractionService(interactionRepository *repository.InteractionRepository, db *pgxpool.Pool, rdb *redis.Client) *InteractionService {
	return &InteractionService{
		interactionRepository: interactionRepository,
		redis:                 rdb,
		db:                    db,
	}
}

func (i InteractionService) PostContent(ctx context.Context, content dto.InteractionRequest, idUser int) error {
	if err := i.interactionRepository.PostContent(ctx, content, idUser); err != nil {
		return nil
	}
	// invalidasi cache
	return nil
}

func (i InteractionService) FollowingUser(ctx context.Context, follow dto.FollowingRequest, idUser int) error {
	if err := i.interactionRepository.FollowingUser(ctx, follow, idUser); err != nil {
		return nil
	}
	// invalidasi cache
	return nil
}

func (i InteractionService) LikePosts(ctx context.Context, like dto.LikeRequest, idUser int) (string, error) {
	tx, err := i.db.Begin(ctx)
	if err != nil {
		log.Println(err)
		return "", err
	}

	data, err := i.interactionRepository.ChecksPosts(ctx, tx, like, idUser)
	if err != nil {
		if err.Error() == "no rows in result set" {
			// msg, err := i.interactionRepository.LikePosts(ctx, tx, like, idUser)
			// if err != nil {
			// 	return "", err
			// }
			// return msg, nil
			log.Println("no rows in result set")
		}
		// return "", err
	}

	defer tx.Rollback(ctx)

	if data != 0 {
		msg, err := i.interactionRepository.DeletePosts(ctx, tx, like, idUser)
		if err != nil {
			return "", err
		}
		if e := tx.Commit(ctx); e != nil {
			log.Println("failed to commit", e.Error())
			return "", e
		}
		return msg, nil
	}

	msg, err := i.interactionRepository.LikePosts(ctx, tx, like, idUser)
	if err != nil {
		return "", err
	}

	if e := tx.Commit(ctx); e != nil {
		log.Println("failed to commit", e.Error())
		return "", e
	}

	return msg, nil
}
