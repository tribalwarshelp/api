package player

import (
	"context"

	"github.com/tribalwarshelp/shared/models"
)

type Repository interface {
	Fetch(ctx context.Context, server string, filter *models.PlayerFilter) ([]*models.Player, int, error)
	FetchPlayerServers(ctx context.Context, playerID ...int) (map[int][]string, error)
}
