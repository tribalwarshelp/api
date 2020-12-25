package player

import (
	"context"

	"github.com/tribalwarshelp/shared/models"
)

type Usecase interface {
	Fetch(ctx context.Context, cfg FetchConfig) ([]*models.Player, int, error)
	GetByID(ctx context.Context, server string, id int) (*models.Player, error)
	SearchPlayer(ctx context.Context, cfg SearchPlayerConfig) ([]*models.FoundPlayer, int, error)
}
