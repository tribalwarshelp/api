package resolvers

import (
	"context"

	"github.com/tribalwarshelp/api/graphql/generated"
	"github.com/tribalwarshelp/api/middleware"
	"github.com/tribalwarshelp/api/utils"
	"github.com/tribalwarshelp/shared/models"
)

func (r *playerResolver) Tribe(ctx context.Context, obj *models.Player) (*models.Tribe, error) {
	if obj.Tribe != nil {
		return obj.Tribe, nil
	}

	return getTribe(ctx, obj.TribeID), nil
}

func (r *playerResolver) Servers(ctx context.Context, obj *models.Player) ([]string, error) {
	langVersionDataLoaders := middleware.LangVersionDataLoadersFromContext(ctx)
	if langVersionDataLoaders != nil {
		serverKey, _ := getServer(ctx)
		if loaders, ok := langVersionDataLoaders[utils.LanguageTagFromServerKey(serverKey)]; ok {
			servers, err := loaders.PlayerServersByID.Load(obj.ID)
			if err == nil {
				return servers, nil
			}
		}
	}
	return []string{}, nil
}

func (r *playerResolver) NameChanges(ctx context.Context, obj *models.Player) ([]*models.PlayerNameChange, error) {
	langVersionDataLoaders := middleware.LangVersionDataLoadersFromContext(ctx)
	if langVersionDataLoaders != nil {
		serverKey, _ := getServer(ctx)
		if loaders, ok := langVersionDataLoaders[utils.LanguageTagFromServerKey(serverKey)]; ok {
			servers, err := loaders.PlayerNameChangesByID.Load(obj.ID)
			if err == nil {
				return servers, nil
			}
		}
	}
	return []*models.PlayerNameChange{}, nil
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
