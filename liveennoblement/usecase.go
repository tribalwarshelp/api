package liveennoblement

import (
	"context"

	"github.com/tribalwarshelp/shared/models"
)

type Usecase interface {
	Fetch(ctx context.Context, server string) ([]*models.LiveEnnoblement, error)
}