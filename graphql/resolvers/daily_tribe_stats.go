package resolvers

import (
	"context"
	"github.com/tribalwarshelp/api/utils"

	"github.com/tribalwarshelp/api/dailytribestats"
	"github.com/tribalwarshelp/api/graphql/generated"
	"github.com/tribalwarshelp/shared/models"
)

func (r *dailyTribeStatsRecordResolver) Tribe(ctx context.Context, obj *models.DailyTribeStats) (*models.Tribe, error) {
	if obj.Tribe != nil {
		return obj.Tribe, nil
	}

	return getTribe(ctx, obj.TribeID), nil
}

func (r *queryResolver) DailyTribeStats(ctx context.Context,
	server string,
	filter *models.DailyTribeStatsFilter,
	limit *int,
	offset *int,
	sort []string) (*generated.DailyTribeStats, error) {
	var err error
	list := &generated.DailyTribeStats{}
	list.Items, list.Total, err = r.DailyTribeStatsUcase.Fetch(ctx, dailytribestats.FetchConfig{
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
