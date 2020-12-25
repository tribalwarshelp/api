package resolvers

import (
	"context"

	"github.com/tribalwarshelp/api/graphql/generated"
	"github.com/tribalwarshelp/api/serverstats"
	"github.com/tribalwarshelp/shared/models"
)

func (r *Resolver) ServerStats(ctx context.Context,
	server string,
	f *models.ServerStatsFilter,
	limit *int,
	offset *int,
	sort []string) (*generated.ServerStats, error) {
	var err error
	list := &generated.ServerStats{}
	list.Items, list.Total, err = r.ServerStatsUcase.Fetch(ctx, serverstats.FetchConfig{
		Server: server,
		Filter: f,
		Sort:   sort,
		Limit:  safeIntPointer(limit, 0),
		Offset: safeIntPointer(offset, 0),
		Count:  shouldCount(ctx),
	})
	return list, err
}
