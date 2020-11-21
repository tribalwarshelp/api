package resolvers

import (
	"context"

	"github.com/tribalwarshelp/api/graphql/generated"
	"github.com/tribalwarshelp/api/middleware"
	"github.com/tribalwarshelp/api/player"
	"github.com/tribalwarshelp/shared/models"
	"github.com/tribalwarshelp/shared/tw"
)

func (r *playerResolver) Tribe(ctx context.Context, obj *models.Player) (*models.Tribe, error) {
	if obj.Tribe != nil {
		return obj.Tribe, nil
	}

	return getTribe(ctx, obj.TribeID), nil
}

func (r *playerResolver) Servers(ctx context.Context, obj *models.Player) ([]string, error) {
	versionDataLoaders := middleware.VersionDataLoadersFromContext(ctx)
	if versionDataLoaders != nil {
		serverKey, _ := getServer(ctx)
		if loaders, ok := versionDataLoaders[tw.VersionCodeFromServerKey(serverKey)]; ok {
			servers, err := loaders.PlayerServersByID.Load(obj.ID)
			if err == nil {
				return servers, nil
			}
		}
	}
	return []string{}, nil
}

func (r *playerResolver) NameChanges(ctx context.Context, obj *models.Player) ([]*models.PlayerNameChange, error) {
	versionDataLoaders := middleware.VersionDataLoadersFromContext(ctx)
	if versionDataLoaders != nil {
		serverKey, _ := getServer(ctx)
		if loaders, ok := versionDataLoaders[tw.VersionCodeFromServerKey(serverKey)]; ok {
			servers, err := loaders.PlayerNameChangesByID.Load(obj.ID)
			if err == nil {
				return servers, nil
			}
		}
	}
	return []*models.PlayerNameChange{}, nil
}

func (r *queryResolver) Players(ctx context.Context,
	server string,
	f *models.PlayerFilter,
	limit *int,
	offset *int,
	sort []string) (*generated.PlayerList, error) {
	defLimit := 0
	defOffset := 0
	if limit == nil {
		limit = &defLimit
	}
	if offset == nil {
		offset = &defOffset
	}

	var err error
	list := &generated.PlayerList{}
	list.Items, list.Total, err = r.PlayerUcase.Fetch(ctx, player.FetchConfig{
		Server: server,
		Filter: f,
		Sort:   sort,
		Limit:  *limit,
		Offset: *offset,
		Count:  true,
	})
	return list, err
}

func (r *queryResolver) Player(ctx context.Context, server string, id int) (*models.Player, error) {
	return r.PlayerUcase.GetByID(ctx, server, id)
}
