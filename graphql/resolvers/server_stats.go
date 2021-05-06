package resolvers

import (
	"context"
	"github.com/Kichiyaki/goutil/safeptr"
	"github.com/tribalwarshelp/shared/tw/twmodel"

	"github.com/tribalwarshelp/api/graphql/generated"
	"github.com/tribalwarshelp/api/serverstats"
)

func (r *Resolver) ServerStats(ctx context.Context,
	server string,
	f *twmodel.ServerStatsFilter,
	limit *int,
	offset *int,
	sort []string) (*generated.ServerStats, error) {
	var err error
	list := &generated.ServerStats{}
	list.Items, list.Total, err = r.ServerStatsUcase.Fetch(ctx, serverstats.FetchConfig{
		Server: server,
		Filter: f,
		Sort:   sort,
		Limit:  safeptr.SafeIntPointer(limit, 0),
		Offset: safeptr.SafeIntPointer(offset, 0),
		Count:  shouldCount(ctx),
		Select: shouldSelectItems(ctx),
	})
	return list, err
}
