package server

import (
	"context"

	"github.com/tribalwarshelp/shared/models"
)

type Repository interface {
	Fetch(ctx context.Context, filter *models.ServerFilter) ([]*models.Server, int, error)
}
