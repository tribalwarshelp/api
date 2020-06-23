package resolvers

import (
	"context"
	"time"

	"github.com/tribalwarshelp/api/graphql/generated"
	"github.com/tribalwarshelp/api/utils"
	"github.com/tribalwarshelp/shared/models"
)

func (r *tribeResolver) CreatedAt(ctx context.Context, obj *models.Tribe) (*time.Time, error) {
	server, _ := getServer(ctx)
	t := formatDate(ctx, utils.LanguageTagFromServerKey(server), obj.CreatedAt)
	return &t, nil
}

func (r *queryResolver) Tribes(ctx context.Context, server string, filter *models.TribeFilter) (*generated.TribesList, error) {
	var err error
	list := &generated.TribesList{}
	list.Items, list.Total, err = r.TribeUcase.Fetch(ctx, server, filter)
	return list, err
}

func (r *queryResolver) Tribe(ctx context.Context, server string, id int) (*models.Tribe, error) {
	return r.TribeUcase.GetByID(ctx, server, id)
}
