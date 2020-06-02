package player

import (
	"context"

	"github.com/tribalwarshelp/shared/models"
)

type Usecase interface {
	Fetch(ctx context.Context, server string, filter *models.PlayerFilter) ([]*models.Player, int, error)
	GetByID(ctx context.Context, server string, id int) (*models.Player, error)
}
