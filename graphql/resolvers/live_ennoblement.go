package resolvers

import (
	"context"

	"github.com/tribalwarshelp/shared/models"
)

func (r *liveEnnoblementResolver) NewOwner(ctx context.Context, obj *models.LiveEnnoblement) (*models.Player, error) {
	if obj.NewOwner != nil {
		return obj.NewOwner, nil
	}

	return getPlayer(ctx, obj.NewOwnerID), nil
}

func (r *liveEnnoblementResolver) OldOwner(ctx context.Context, obj *models.LiveEnnoblement) (*models.Player, error) {
	if obj.OldOwner != nil {
		return obj.OldOwner, nil
	}

	return getPlayer(ctx, obj.OldOwnerID), nil
}

func (r *liveEnnoblementResolver) Village(ctx context.Context, obj *models.LiveEnnoblement) (*models.Village, error) {
	if obj.Village != nil {
		return obj.Village, nil
	}

	return getVillage(ctx, obj.VillageID), nil
}

func (r *queryResolver) LiveEnnoblements(ctx context.Context, server string) ([]*models.LiveEnnoblement, error) {
	return r.EnnoblementUcase.FetchLiveEnnoblements(ctx, server)
}
