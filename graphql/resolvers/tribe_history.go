package resolvers

import (
	"context"
	"github.com/tribalwarshelp/api/utils"

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
	var err error
	list := &generated.TribeHistory{}
	list.Items, list.Total, err = r.TribeHistoryUcase.Fetch(ctx, tribehistory.FetchConfig{
		Filter: f,
		Sort:   sort,
		Limit:  utils.SafeIntPointer(limit, 0),
		Offset: utils.SafeIntPointer(offset, 0),
		Count:  shouldCount(ctx),
		Server: server,
	})
	return list, err
}
