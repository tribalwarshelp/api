package resolvers

import (
	"context"

	"github.com/tribalwarshelp/api/graphql/generated"
	"github.com/tribalwarshelp/api/tribehistory"
	"github.com/tribalwarshelp/shared/models"
)

func (r *tribeHistoryRecordResolver) Tribe(ctx context.Context, obj *models.TribeHistory) (*models.Tribe, error) {
	if obj.Tribe != nil {
		return obj.Tribe, nil
	}

	return getTribe(ctx, obj.TribeID), nil
}

func (r *Resolver) TribeHistory(ctx context.Context,
	server string,
	f *models.TribeHistoryFilter,
	limit *int,
	offset *int,
	sort []string) (*generated.TribeHistory, error) {
	defLimit := 0
	defOffset := 0
	if limit == nil {
		limit = &defLimit
	}
	if offset == nil {
		offset = &defOffset
	}

	var err error
	list := &generated.TribeHistory{}
	list.Items, list.Total, err = r.TribeHistoryUcase.Fetch(ctx, tribehistory.FetchConfig{
		Filter: f,
		Sort:   sort,
		Limit:  *limit,
		Offset: *offset,
		Count:  true,
		Server: server,
	})
	return list, err
}
