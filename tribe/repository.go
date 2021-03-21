package tribe

import (
	"context"

	"github.com/tribalwarshelp/shared/models"
)

type FetchConfig struct {
	Server string
	Filter *models.TribeFilter
	Count  bool
	Select bool
	Sort   []string
	Limit  int
	Offset int
}

type SearchTribeConfig struct {
	Version string
	Query   string
	Count   bool
	Sort    []string
	Limit   int
	Offset  int
}

type Repository interface {
	Fetch(ctx context.Context, cfg FetchConfig) ([]*models.Tribe, int, error)
	SearchTribe(ctx context.Context, cfg SearchTribeConfig) ([]*models.FoundTribe, int, error)
}
