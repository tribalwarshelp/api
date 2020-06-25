package resolvers

import (
	"context"

	"github.com/tribalwarshelp/api/graphql/generated"
	"github.com/tribalwarshelp/shared/models"
)

func (r *dailyPlayerStatsRecordResolver) Player(ctx context.Context, obj *models.DailyPlayerStats) (*models.Player, error) {
	if obj.Player != nil {
		return obj.Player, nil
	}

	return getPlayer(ctx, obj.PlayerID), nil
}

func (r *queryResolver) DailyPlayerStats(ctx context.Context, server string, filter *models.DailyPlayerStatsFilter) (*generated.DailyPlayerStats, error) {
	var err error
	list := &generated.DailyPlayerStats{}
	list.Items, list.Total, err = r.DailyPlayerStatsUcase.Fetch(ctx, server, filter)
	return list, err
}
