package service

import (
	"context"

	"github.com/ariboss89/coffee-morning-services/internal/dto"
	"github.com/ariboss89/coffee-morning-services/internal/repository"
)

type InteractionService struct {
	interactionRepository *repository.InteractionRepository
}

func NewInteractionService(interactionRepository *repository.InteractionRepository) *InteractionService {
	return &InteractionService{
		interactionRepository: interactionRepository,
	}
}

func (i InteractionService) PostContent(ctx context.Context, content dto.InteractionRequest, idUser int) error {
	if err := i.interactionRepository.PostContent(ctx, content, idUser); err != nil {
		return nil
	}
	// invalidasi cache
	return nil
}
