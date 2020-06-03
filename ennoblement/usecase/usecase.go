package usecase

import (
	"context"

	"github.com/tribalwarshelp/api/ennoblement"
	"github.com/tribalwarshelp/shared/models"
)

type usecase struct {
	repo ennoblement.Repository
}

func New(repo ennoblement.Repository) ennoblement.Usecase {
	return &usecase{
		repo,
	}
}

func (ucase *usecase) Fetch(ctx context.Context, server string) ([]*models.Ennoblement, error) {
	return ucase.repo.Fetch(ctx, server)
}
