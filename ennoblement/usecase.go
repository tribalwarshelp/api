package ennoblement

import (
	"context"

	"github.com/tribalwarshelp/shared/models"
)

type Usecase interface {
	Fetch(ctx context.Context, cfg FetchConfig) ([]*models.Ennoblement, int, error)
	FetchLiveEnnoblements(ctx context.Context, server string) ([]*models.LiveEnnoblement, error)
}
