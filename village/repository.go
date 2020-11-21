package village

import (
	"context"

	"github.com/tribalwarshelp/shared/models"
)

type FetchConfig struct {
	Server  string
	Filter  *models.VillageFilter
	Columns []string
	Count   bool
	Sort    []string
	Limit   int
	Offset  int
}

type Repository interface {
	Fetch(ctx context.Context, cfg FetchConfig) ([]*models.Village, int, error)
}
