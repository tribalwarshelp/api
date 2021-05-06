package resolvers

import (
	"context"
	"github.com/Kichiyaki/goutil/safeptr"
	"github.com/tribalwarshelp/shared/tw/twmodel"

	"github.com/tribalwarshelp/api/graphql/generated"
	"github.com/tribalwarshelp/api/village"
)

func (r *villageResolver) Player(ctx context.Context, obj *twmodel.Village) (*twmodel.Player, error) {
	if obj.Player != nil {
		return obj.Player, nil
	}

	return getPlayer(ctx, obj.PlayerID), nil
}

func (r *queryResolver) Villages(ctx context.Context,
	server string,
	f *twmodel.VillageFilter,
	limit *int,
	offset *int,
	sort []string) (*generated.VillageList, error) {
	var err error
	list := &generated.VillageList{}
	list.Items, list.Total, err = r.VillageUcase.Fetch(ctx, village.FetchConfig{
		Filter: f,
		Sort:   sort,
		Limit:  safeptr.SafeIntPointer(limit, 0),
		Offset: safeptr.SafeIntPointer(offset, 0),
		Select: shouldSelectItems(ctx),
		Count:  shouldCount(ctx),
		Server: server,
	})
	return list, err
}

func (r *queryResolver) Village(ctx context.Context, server string, id int) (*twmodel.Village, error) {
	return r.VillageUcase.GetByID(ctx, server, id)
}
