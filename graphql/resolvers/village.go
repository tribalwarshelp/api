package resolvers

import (
	"context"
	"github.com/tribalwarshelp/api/utils"

	"github.com/tribalwarshelp/api/graphql/generated"
	"github.com/tribalwarshelp/api/village"
	"github.com/tribalwarshelp/shared/models"
)

func (r *villageResolver) Player(ctx context.Context, obj *models.Village) (*models.Player, error) {
	if obj.Player != nil {
		return obj.Player, nil
	}

	return getPlayer(ctx, obj.PlayerID), nil
}

func (r *queryResolver) Villages(ctx context.Context,
	server string,
	f *models.VillageFilter,
	limit *int,
	offset *int,
	sort []string) (*generated.VillageList, error) {
	var err error
	list := &generated.VillageList{}
	list.Items, list.Total, err = r.VillageUcase.Fetch(ctx, village.FetchConfig{
		Filter: f,
		Sort:   sort,
		Limit:  utils.SafeIntPointer(limit, 0),
		Offset: utils.SafeIntPointer(offset, 0),
		Count:  shouldCount(ctx),
		Server: server,
	})
	return list, err
}

func (r *queryResolver) Village(ctx context.Context, server string, id int) (*models.Village, error) {
	return r.VillageUcase.GetByID(ctx, server, id)
}
