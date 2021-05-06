package resolvers

import (
	"context"
	"github.com/Kichiyaki/goutil/safeptr"
	"github.com/tribalwarshelp/shared/tw/twmodel"

	"github.com/tribalwarshelp/api/dailyplayerstats"
	"github.com/tribalwarshelp/api/graphql/generated"
)

func (r *dailyPlayerStatsRecordResolver) Player(ctx context.Context, obj *twmodel.DailyPlayerStats) (*twmodel.Player, error) {
	if obj.Player != nil {
		return obj.Player, nil
	}

	return getPlayer(ctx, obj.PlayerID), nil
}

func (r *queryResolver) DailyPlayerStats(ctx context.Context,
	server string,
	filter *twmodel.DailyPlayerStatsFilter,
	limit *int,
	offset *int,
	sort []string) (*generated.DailyPlayerStats, error) {
	var err error
	list := &generated.DailyPlayerStats{}
	list.Items, list.Total, err = r.DailyPlayerStatsUcase.Fetch(ctx, dailyplayerstats.FetchConfig{
		Server: server,
		Filter: filter,
		Sort:   sort,
		Limit:  safeptr.SafeIntPointer(limit, 0),
		Offset: safeptr.SafeIntPointer(offset, 0),
		Select: shouldSelectItems(ctx),
		Count:  shouldCount(ctx),
	})
	return list, err
}
