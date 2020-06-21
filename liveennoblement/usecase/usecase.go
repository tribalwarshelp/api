package usecase

import (
	"context"

	"github.com/tribalwarshelp/api/liveennoblement"
	"github.com/tribalwarshelp/shared/models"
)

type usecase struct {
	repo liveennoblement.Repository
}

func New(repo liveennoblement.Repository) liveennoblement.Usecase {
	return &usecase{
		repo,
	}
}

func (ucase *usecase) Fetch(ctx context.Context, server string) ([]*models.LiveEnnoblement, error) {
	return ucase.repo.Fetch(ctx, server)
}
