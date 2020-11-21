package tribechange

import (
	"context"

	"github.com/tribalwarshelp/shared/models"
)

type FetchConfig struct {
	Server string
	Filter *models.TribeChangeFilter
	Count  bool
	Sort   []string
	Limit  int
	Offset int
}

type Repository interface {
	Fetch(ctx context.Context, cfg FetchConfig) ([]*models.TribeChange, int, error)
}
