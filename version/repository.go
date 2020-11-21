package version

import (
	"context"

	"github.com/tribalwarshelp/shared/models"
)

type FetchConfig struct {
	Filter *models.VersionFilter
	Count  bool
	Sort   []string
	Limit  int
	Offset int
}

type Repository interface {
	Fetch(ctx context.Context, cfg FetchConfig) ([]*models.Version, int, error)
}
