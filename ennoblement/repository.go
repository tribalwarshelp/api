package ennoblement

import (
	"context"

	"github.com/tribalwarshelp/shared/models"
)

type Repository interface {
	Fetch(ctx context.Context, server string) ([]*models.Ennoblement, error)
}
