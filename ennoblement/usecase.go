package ennoblement

import (
	"context"

	"github.com/tribalwarshelp/shared/models"
)

type Usecase interface {
	Fetch(ctx context.Context, server string) ([]*models.Ennoblement, error)
}
