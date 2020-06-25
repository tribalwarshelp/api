package server

import (
	"context"

	"github.com/tribalwarshelp/shared/models"
)

type FetchConfig struct {
	Filter *models.ServerFilter
	Count  bool
}

type Repository interface {
	Fetch(ctx context.Context, cfg FetchConfig) ([]*models.Server, int, error)
}
