package ennoblement

import (
	"context"

	"github.com/tribalwarshelp/shared/models"
)

type FetchConfig struct {
	Server string
	Filter *models.EnnoblementFilter
	Count  bool
	Sort   []string
	Limit  int
	Offset int
}

type Repository interface {
	Fetch(ctx context.Context, cfg FetchConfig) ([]*models.Ennoblement, int, error)
}
