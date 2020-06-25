package tribe

import (
	"context"

	"github.com/tribalwarshelp/shared/models"
)

type FetchConfig struct {
	Server string
	Filter *models.TribeFilter
	Count  bool
}

type Repository interface {
	Fetch(ctx context.Context, cfg FetchConfig) ([]*models.Tribe, int, error)
}
