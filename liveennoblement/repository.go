package liveennoblement

import (
	"context"

	"github.com/tribalwarshelp/shared/models"
)

type Repository interface {
	Fetch(ctx context.Context, server string) ([]*models.LiveEnnoblement, error)
}
