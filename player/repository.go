package player

import (
	"context"

	"github.com/tribalwarshelp/shared/models"
)

type FetchConfig struct {
	Server string
	Filter *models.PlayerFilter
	Count  bool
	Sort   []string
	Limit  int
	Offset int
}

type Repository interface {
	Fetch(ctx context.Context, cfg FetchConfig) ([]*models.Player, int, error)
	FetchNameChanges(ctx context.Context, code models.VersionCode, playerID ...int) (map[int][]*models.PlayerNameChange, error)
	FetchPlayerServers(ctx context.Context, code models.VersionCode, playerID ...int) (map[int][]string, error)
}
