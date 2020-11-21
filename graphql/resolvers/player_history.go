package resolvers

import (
	"context"

	"github.com/tribalwarshelp/api/graphql/generated"
	"github.com/tribalwarshelp/api/playerhistory"
	"github.com/tribalwarshelp/shared/models"
)

func (r *playerHistoryRecordResolver) Player(ctx context.Context, obj *models.PlayerHistory) (*models.Player, error) {
	if obj.Player != nil {
		return obj.Player, nil
	}

	return getPlayer(ctx, obj.PlayerID), nil
}

func (r *playerHistoryRecordResolver) Tribe(ctx context.Context, obj *models.PlayerHistory) (*models.Tribe, error) {
	if obj.Tribe != nil {
		return obj.Tribe, nil
	}

	return getTribe(ctx, obj.TribeID), nil
}

func (r *Resolver) PlayerHistory(ctx context.Context,
	server string,
	f *models.PlayerHistoryFilter,
	limit *int,
	offset *int,
	sort []string) (*generated.PlayerHistory, error) {
	defLimit := 0
	defOffset := 0
	if limit == nil {
		limit = &defLimit
	}
	if offset == nil {
		offset = &defOffset
	}

	var err error
	list := &generated.PlayerHistory{}
	list.Items, list.Total, err = r.PlayerHistoryUcase.Fetch(ctx, playerhistory.FetchConfig{
		Server: server,
		Filter: f,
		Sort:   sort,
		Limit:  *limit,
		Offset: *offset,
		Count:  true,
	})
	return list, err
}
