package resolvers

import (
	"context"
	"github.com/Kichiyaki/goutil/safeptr"
	"github.com/tribalwarshelp/shared/tw/twmodel"

	"github.com/tribalwarshelp/api/graphql/generated"
	"github.com/tribalwarshelp/api/middleware"
	"github.com/tribalwarshelp/api/player"
)

func (r *playerResolver) Tribe(ctx context.Context, obj *twmodel.Player) (*twmodel.Tribe, error) {
	if obj.Tribe != nil {
		return obj.Tribe, nil
	}

	return getTribe(ctx, obj.TribeID), nil
}

func (r *playerResolver) Servers(ctx context.Context, obj *twmodel.Player) ([]string, error) {
	versionDataLoaders := middleware.VersionDataLoadersFromContext(ctx)
	if versionDataLoaders != nil {
		serverKey, _ := getServer(ctx)
		if loaders, ok := versionDataLoaders[twmodel.VersionCodeFromServerKey(serverKey)]; ok {
			servers, err := loaders.PlayerServersByID.Load(obj.ID)
			if err == nil {
				return servers, nil
			}
		}
	}
	return []string{}, nil
}

func (r *playerResolver) NameChanges(ctx context.Context, obj *twmodel.Player) ([]*twmodel.PlayerNameChange, error) {
	versionDataLoaders := middleware.VersionDataLoadersFromContext(ctx)
	if versionDataLoaders != nil {
		serverKey, _ := getServer(ctx)
		if loaders, ok := versionDataLoaders[twmodel.VersionCodeFromServerKey(serverKey)]; ok {
			servers, err := loaders.PlayerNameChangesByID.Load(obj.ID)
			if err == nil {
				return servers, nil
			}
		}
	}
	return []*twmodel.PlayerNameChange{}, nil
}

func (r *queryResolver) Players(ctx context.Context,
	server string,
	f *twmodel.PlayerFilter,
	limit *int,
	offset *int,
	sort []string) (*generated.PlayerList, error) {
	var err error
	list := &generated.PlayerList{}
	list.Items, list.Total, err = r.PlayerUcase.Fetch(ctx, player.FetchConfig{
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

func (r *queryResolver) Player(ctx context.Context, server string, id int) (*twmodel.Player, error) {
	return r.PlayerUcase.GetByID(ctx, server, id)
}

func (r *queryResolver) SearchPlayer(ctx context.Context,
	version string,
	name *string,
	id *int,
	limit *int,
	offset *int,
	sort []string) (*generated.FoundPlayerList, error) {
	var err error
	list := &generated.FoundPlayerList{}
	list.Items, list.Total, err = r.PlayerUcase.SearchPlayer(ctx, player.SearchPlayerConfig{
		Sort:    sort,
		Limit:   safeptr.SafeIntPointer(limit, 0),
		Offset:  safeptr.SafeIntPointer(offset, 0),
		Version: version,
		Name:    safeptr.SafeStringPointer(name, ""),
		ID:      safeptr.SafeIntPointer(id, 0),
		Count:   shouldCount(ctx),
	})
	return list, err
}
