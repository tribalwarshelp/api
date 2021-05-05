package resolvers

import (
	"context"
	"github.com/Kichiyaki/goutil/safeptr"
	"github.com/tribalwarshelp/shared/tw/twmodel"

	"github.com/tribalwarshelp/api/dailytribestats"
	"github.com/tribalwarshelp/api/graphql/generated"
)

func (r *dailyTribeStatsRecordResolver) Tribe(ctx context.Context, obj *twmodel.DailyTribeStats) (*twmodel.Tribe, error) {
	if obj.Tribe != nil {
		return obj.Tribe, nil
	}

	return getTribe(ctx, obj.TribeID), nil
}

func (r *queryResolver) DailyTribeStats(ctx context.Context,
	server string,
	filter *twmodel.DailyTribeStatsFilter,
	limit *int,
	offset *int,
	sort []string) (*generated.DailyTribeStats, error) {
	var err error
	list := &generated.DailyTribeStats{}
	list.Items, list.Total, err = r.DailyTribeStatsUcase.Fetch(ctx, dailytribestats.FetchConfig{
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
