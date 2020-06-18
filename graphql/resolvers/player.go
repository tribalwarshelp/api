package resolvers

import (
	"context"

	"github.com/99designs/gqlgen/graphql"
	"github.com/tribalwarshelp/api/graphql/generated"
	"github.com/tribalwarshelp/api/middleware"
	"github.com/tribalwarshelp/shared/models"
)

func (r *playerResolver) Tribe(ctx context.Context, obj *models.Player) (*models.Tribe, error) {
	if obj.Tribe != nil {
		return obj.Tribe, nil
	}

	if server, ok := getServer(graphql.GetFieldContext(ctx)); ok {
		dataloaders := middleware.ServerDataLoadersFromContext(ctx)
		if dataloaders != nil {
			if dataloader, ok := dataloaders[server]; ok {
				tribe, _ := dataloader.TribeByID.Load(obj.TribeID)
				if tribe != nil {
					return tribe, nil
				}
			}
		}
	}

	return nil, nil
}

func (r *queryResolver) Players(ctx context.Context, server string, filter *models.PlayerFilter) (*generated.PlayersList, error) {
	var err error
	list := &generated.PlayersList{}
	list.Items, list.Total, err = r.PlayerUcase.Fetch(ctx, server, filter)
	return list, err
}

func (r *queryResolver) Player(ctx context.Context, server string, id int) (*models.Player, error) {
	return r.PlayerUcase.GetByID(ctx, server, id)
}
