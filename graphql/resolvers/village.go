package resolvers

import (
	"context"

	"github.com/tribalwarshelp/api/graphql/generated"
	"github.com/tribalwarshelp/shared/models"
)

func (r *villageResolver) Player(ctx context.Context, obj *models.Village) (*models.Player, error) {
	if obj.Player != nil {
		return obj.Player, nil
	}

	return getPlayer(ctx, obj.PlayerID), nil
}

func (r *queryResolver) Villages(ctx context.Context, server string, filter *models.VillageFilter) (*generated.VillagesList, error) {
	var err error
	list := &generated.VillagesList{}
	list.Items, list.Total, err = r.VillageUcase.Fetch(ctx, server, filter)
	return list, err
}

func (r *queryResolver) Village(ctx context.Context, server string, id int) (*models.Village, error) {
	return r.VillageUcase.GetByID(ctx, server, id)
}
