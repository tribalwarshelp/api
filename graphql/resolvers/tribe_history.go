package resolvers

import (
	"context"

	"github.com/tribalwarshelp/api/graphql/generated"
	"github.com/tribalwarshelp/shared/models"
)

func (r *tribeHistoryRecordResolver) Tribe(ctx context.Context, obj *models.TribeHistory) (*models.Tribe, error) {
	if obj.Tribe != nil {
		return obj.Tribe, nil
	}

	return getTribe(ctx, obj.TribeID), nil
}

func (r *Resolver) TribeHistory(ctx context.Context, server string, filter *models.TribeHistoryFilter) (*generated.TribeHistory, error) {
	var err error
	list := &generated.TribeHistory{}
	list.Items, list.Total, err = r.TribeHistoryUcase.Fetch(ctx, server, filter)
	return list, err
}
