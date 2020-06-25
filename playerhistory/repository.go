package playerhistory

import (
	"context"

	"github.com/tribalwarshelp/shared/models"
)

type FetchConfig struct {
	Server string
	Filter *models.PlayerHistoryFilter
	Count  bool
}

type Repository interface {
	Fetch(ctx context.Context, cfg FetchConfig) ([]*models.PlayerHistory, int, error)
}
