package resolvers

import (
	"context"

	"github.com/99designs/gqlgen/graphql"
	"github.com/tribalwarshelp/api/graphql/generated"
	"github.com/tribalwarshelp/api/middleware"
	"github.com/tribalwarshelp/shared/models"
)

func (r *villageResolver) Player(ctx context.Context, obj *models.Village) (*models.Player, error) {
	if obj.Player != nil {
		return obj.Player, nil
	}

	rctx := graphql.GetFieldContext(ctx)
	server, ok := rctx.Parent.Parent.Parent.Args["server"].(string)
	if ok {
		dataloaders := middleware.DataLoadersFromContext(ctx)
		if dataloaders != nil {
			if dataloader, ok := dataloaders[server]; ok {
				tribe, _ := dataloader.PlayerByID.Load(obj.PlayerID)
				if tribe != nil {
					return tribe, nil
				}
			}
		}
	}

	return nil, nil
}

func (r *queryResolver) Villages(ctx context.Context, server string, filter *models.VillageFilter) (*generated.VillageList, error) {
	var err error
	list := &generated.VillageList{}
	list.Items, list.Total, err = r.VillageUcase.Fetch(ctx, server, filter)
	return list, err
}

func (r *queryResolver) Village(ctx context.Context, server string, id int) (*models.Village, error) {
	return r.VillageUcase.GetByID(ctx, server, id)
}
