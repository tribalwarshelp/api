package ennoblement

import (
	"context"

	"github.com/tribalwarshelp/shared/models"
)

type Repository interface {
	Fetch(ctx context.Context, server string, filter *models.EnnoblementFilter) ([]*models.Ennoblement, int, error)
}
