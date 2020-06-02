package server

import (
	"context"

	"github.com/tribalwarshelp/shared/models"
)

type Usecase interface {
	Fetch(ctx context.Context, filter *models.ServerFilter) ([]*models.Server, int, error)
	GetByKey(ctx context.Context, key string) (*models.Server, error)
}
