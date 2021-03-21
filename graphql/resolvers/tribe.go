package resolvers

import (
	"context"
	"github.com/tribalwarshelp/api/utils"

	"github.com/tribalwarshelp/api/graphql/generated"
	"github.com/tribalwarshelp/api/tribe"
	"github.com/tribalwarshelp/shared/models"
)

func (r *queryResolver) Tribes(ctx context.Context,
	server string,
	f *models.TribeFilter,
	limit *int,
	offset *int,
	sort []string) (*generated.TribeList, error) {
	var err error
	list := &generated.TribeList{}
	list.Items, list.Total, err = r.TribeUcase.Fetch(ctx, tribe.FetchConfig{
		Server: server,
		Filter: f,
		Sort:   sort,
		Limit:  utils.SafeIntPointer(limit, 0),
		Offset: utils.SafeIntPointer(offset, 0),
		Count:  shouldCount(ctx),
		Select: shouldSelectItems(ctx),
	})
	return list, err
}

func (r *queryResolver) Tribe(ctx context.Context, server string, id int) (*models.Tribe, error) {
	return r.TribeUcase.GetByID(ctx, server, id)
}

func (r *queryResolver) SearchTribe(ctx context.Context,
	version string,
	query string,
	limit *int,
	offset *int,
	sort []string) (*generated.FoundTribeList, error) {
	var err error
	list := &generated.FoundTribeList{}
	list.Items, list.Total, err = r.TribeUcase.SearchTribe(ctx, tribe.SearchTribeConfig{
		Sort:    sort,
		Limit:   utils.SafeIntPointer(limit, 0),
		Offset:  utils.SafeIntPointer(offset, 0),
		Version: version,
		Query:   query,
		Count:   shouldCount(ctx),
	})
	return list, err
}
