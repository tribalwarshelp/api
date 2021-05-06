package resolvers

import (
	"context"
	"github.com/Kichiyaki/goutil/safeptr"
	"github.com/tribalwarshelp/shared/tw/twmodel"

	"github.com/tribalwarshelp/api/graphql/generated"
	"github.com/tribalwarshelp/api/playerhistory"
)

func (r *playerHistoryRecordResolver) Player(ctx context.Context, obj *twmodel.PlayerHistory) (*twmodel.Player, error) {
	if obj.Player != nil {
		return obj.Player, nil
	}

	return getPlayer(ctx, obj.PlayerID), nil
}

func (r *playerHistoryRecordResolver) Tribe(ctx context.Context, obj *twmodel.PlayerHistory) (*twmodel.Tribe, error) {
	if obj.Tribe != nil {
		return obj.Tribe, nil
	}

	return getTribe(ctx, obj.TribeID), nil
}

func (r *Resolver) PlayerHistory(ctx context.Context,
	server string,
	f *twmodel.PlayerHistoryFilter,
	limit *int,
	offset *int,
	sort []string) (*generated.PlayerHistory, error) {
	var err error
	list := &generated.PlayerHistory{}
	list.Items, list.Total, err = r.PlayerHistoryUcase.Fetch(ctx, playerhistory.FetchConfig{
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
