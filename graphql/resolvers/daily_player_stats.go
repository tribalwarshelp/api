package resolvers

import (
	"context"
	"github.com/tribalwarshelp/api/utils"

	"github.com/tribalwarshelp/api/dailyplayerstats"
	"github.com/tribalwarshelp/api/graphql/generated"
	"github.com/tribalwarshelp/shared/models"
)

func (r *dailyPlayerStatsRecordResolver) Player(ctx context.Context, obj *models.DailyPlayerStats) (*models.Player, error) {
	if obj.Player != nil {
		return obj.Player, nil
	}

	return getPlayer(ctx, obj.PlayerID), nil
}

func (r *queryResolver) DailyPlayerStats(ctx context.Context,
	server string,
	filter *models.DailyPlayerStatsFilter,
	limit *int,
	offset *int,
	sort []string) (*generated.DailyPlayerStats, error) {
	var err error
	list := &generated.DailyPlayerStats{}
	list.Items, list.Total, err = r.DailyPlayerStatsUcase.Fetch(ctx, dailyplayerstats.FetchConfig{
		Server: server,
		Filter: filter,
		Sort:   sort,
		Limit:  utils.SafeIntPointer(limit, 0),
		Offset: utils.SafeIntPointer(offset, 0),
		Select: shouldSelectItems(ctx),
		Count:  shouldCount(ctx),
	})
	return list, err
}
