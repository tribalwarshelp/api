package tribe

import (
	"context"
	"github.com/tribalwarshelp/shared/tw/twmodel"
)

type Usecase interface {
	Fetch(ctx context.Context, cfg FetchConfig) ([]*twmodel.Tribe, int, error)
	GetByID(ctx context.Context, server string, id int) (*twmodel.Tribe, error)
	SearchTribe(ctx context.Context, cfg SearchTribeConfig) ([]*twmodel.FoundTribe, int, error)
}
