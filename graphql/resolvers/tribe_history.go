package resolvers

import (
	"context"
	"github.com/Kichiyaki/goutil/safeptr"
	"github.com/tribalwarshelp/shared/tw/twmodel"

	"github.com/tribalwarshelp/api/graphql/generated"
	"github.com/tribalwarshelp/api/tribehistory"
)

func (r *tribeHistoryRecordResolver) Tribe(ctx context.Context, obj *twmodel.TribeHistory) (*twmodel.Tribe, error) {
	if obj.Tribe != nil {
		return obj.Tribe, nil
	}

	return getTribe(ctx, obj.TribeID), nil
}

func (r *Resolver) TribeHistory(ctx context.Context,
	server string,
	f *twmodel.TribeHistoryFilter,
	limit *int,
	offset *int,
	sort []string) (*generated.TribeHistory, error) {
	var err error
	list := &generated.TribeHistory{}
	list.Items, list.Total, err = r.TribeHistoryUcase.Fetch(ctx, tribehistory.FetchConfig{
		Filter: f,
		Sort:   sort,
		Limit:  safeptr.SafeIntPointer(limit, 0),
		Offset: safeptr.SafeIntPointer(offset, 0),
		Count:  shouldCount(ctx),
		Select: shouldSelectItems(ctx),
		Server: server,
	})
	return list, err
}
